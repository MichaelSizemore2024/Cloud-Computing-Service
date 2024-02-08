package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/gocql/gocql" // Scylla Drivers

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	Routes "dbmanager/common" // Import the generated code from protofiles
)

// Variables that can be passed in via command line such as -port=50051
var (
	port  = flag.Int("port", 50051, "The server port")
	debug = flag.Bool("debug", false, "Debug output")
)

// Server implements the interface
type server struct {
	Routes.UnimplementedDBGenericServer // Add this line to embed the unimplemented methods
	mu                                  sync.Mutex
	cluster                             *gocql.ClusterConfig
}

// Newserver initializes and returns a new server.
func NewServer() *server {
	return &server{
		cluster: initDatabaseCluster(),
	}
}

// initDatabaseCluster initializes the Cassandra/ScyllaDB cluster configuration.
func initDatabaseCluster() *gocql.ClusterConfig {
	cluster := gocql.NewCluster("scylla") // Add ScyllaDB node IP or name here (scylla in this instance)
	cluster.Consistency = gocql.One       // Set the consistency level
	return cluster
}

/* CRUD operations
* Insert
* Selelct
* Update
* Delete
* Drop Table
 */

/*
* Handles INSERT requests and stores the provided protobuf messages in the database.
*
 * This method updates`` the table based on the provided condition
 *
 * Parameters:
 *   - ctx: A context object for the request
 *   - request: A protobuf containing information about the keyspace and table to drop
 *
 * Returns:
 *   - *ProtobufErrorResponse: A response containing error information if an error occurs.
*/
func (s *server) Insert(ctx context.Context, request *Routes.ProtobufInsertRequest) (*Routes.ProtobufErrorResponse, error) {
	// Initialize debug counter
	counter := 0

	// Loops over all protobufs received and process them
	for _, any := range request.Protobufs {

		// Unmarshal each Any message to a ProtoMessage
		m, err := any.UnmarshalNew()

		// Return error to client in a list of strings if encountered
		if err != nil {
			return &Routes.ProtobufErrorResponse{Errs: []string{err.Error()}}, err
		}

		msg_desc := m.ProtoReflect().Descriptor() // find the message descriptor
		fields := msg_desc.Fields()               // get all the fields in the proto message

		queryCols := []string{} // list of strings to use strings.Join later for our column names to insert into table
		queryQms := []string{}  // just to keep track of ?'s in the query to add
		serializedAny, err := proto.Marshal(m)
		values := []interface{}{serializedAny}                   // all the values of each field, input for a variadic function/var, timeUUID for uniqie primary key for now
		queryColTypes := []string{`serial_msg BLOB PRIMARY KEY`} // for CREATE TABLE, generate the id with the curr datetime to make it unique

		// Adds the serialized message to the table
		queryCols = append(queryCols, "serial_msg")
		queryQms = append(queryQms, "?")

		// Loop over each field and infer data
		for i := 0; i < fields.Len(); i++ {
			curField := fields.Get(i)

			queryCols = append(queryCols, string(curField.FullName().Name())) // append the field name to our col list
			queryQms = append(queryQms, "?")                                  // append ?'s for values portion of cql cmds

			convertedType, convertedValue := convertKindAndValueToCQL(curField.Kind(), m.ProtoReflect().Get(curField))
			values = append(values, convertedValue)                                                     // get the value of the field
			queryColTypes = append(queryColTypes, fmt.Sprintf("%s %s", curField.Name(), convertedType)) // NAME TYPE for CQL create table
		}

		// Lock to ensure only one goroutine is executed at a time
		s.mu.Lock()
		defer s.mu.Unlock()

		// Create a session to interact with the database
		session, err := s.cluster.CreateSession()
		if err != nil {
			log.Printf("Error opening cluster session")
			return nil, err
		}
		defer session.Close()

		// Put together insert query
		tableName := msg_desc.Name()
		ks := request.Keyspace

		// Create the keyspace if it does not exist, note ` vs " matters
		// TODO: CHANGE CLASS AND REPLICATION FACTOR TO ?
		ksQuery := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy','replication_factor': 3}`, ks)
		if ksQueryErr := session.Query(ksQuery).Exec(); ksQueryErr != nil {
			log.Printf("Error creating keyspace: %s\n query: %s", ksQueryErr, ksQuery)
			return &Routes.ProtobufErrorResponse{Errs: []string{ksQueryErr.Error()}}, ksQueryErr
		}

		// Create a table if it does not exist
		// sometimes cql could try to insert a reserved keyword as a column name and it wont work?
		tableQuery := fmt.Sprint(`CREATE TABLE IF NOT EXISTS `, ks, `.`, tableName, ` (`, strings.Join(queryColTypes, `, `), `)`)
		if tableQueryErr := session.Query(tableQuery).Exec(); tableQueryErr != nil {
			log.Printf("Error creating table: %s\n query: %s", tableQueryErr, tableQuery)
			return &Routes.ProtobufErrorResponse{Errs: []string{tableQueryErr.Error()}}, tableQueryErr
		}

		// handle dynamic insert query
		insertQuery := fmt.Sprint(`INSERT INTO `, ks, `.`, tableName, ` ( `, strings.Join(queryCols, `, `), `) VALUES (`, strings.Join(queryQms, `, `), `)`) // strings.Join is a nice method
		if insertQueryErr := session.Query(insertQuery, values...).Exec(); insertQueryErr != nil {                                                           // look up Variadic Functions for more on the ... syntax
			log.Printf("Error inserting data: %s\n query: %s", insertQueryErr, insertQuery)
			return &Routes.ProtobufErrorResponse{Errs: []string{insertQueryErr.Error()}}, insertQueryErr
		}

		// Increment counter
		counter++
	}

	// Prints out total # of rows updated
	if *debug {
		fmt.Println("Inserted", counter, "entries")
	}

	return &Routes.ProtobufErrorResponse{Errs: []string{}}, nil
}

/*
* Select handles SELECT requests based on specified conditions.

* Parameters:
*   - ctx: A context object for the request
*   - request: A protobuf containing information about the keyspace, table, and constraint
*
* Returns:
*   - *ProtobufErrorResponse: A response containing error information if an error occurs.
 */
func (s *server) Select(ctx context.Context, request *Routes.ProtobufSelectRequest) (*Routes.ProtobufSelectResponse, error) {
	// Lock to ensure only one goroutine is executed at a time
	s.mu.Lock()
	defer s.mu.Unlock()

	// Retrieves data from received protobuf
	tableName := request.Table
	ks := request.Keyspace
	column := request.Column
	constraint := request.Constraint

	// Create a session to interact with the database
	session, err := s.cluster.CreateSession()
	if err != nil {
		log.Printf("Error opening cluster session")
		return nil, err
	}
	defer session.Close()

	// Create index so we can search through constraints (otherwise we can only search via PK)
	if err := createIndex(session, ks, tableName, column); err != nil {
		return &Routes.ProtobufSelectResponse{
			Response:  err.Error(), // Use the error message as the response
			Protobufs: nil,         // Initialize the protobufs slice
		}, err
	}

	// Gets the value type of the column being used
	// Might remove this if we change how the condition is passed in
	columnType, err := getColumnType(session, ks, tableName, column)
	if err != nil {
		log.Printf("Error getting column type: %s", err)
		return &Routes.ProtobufSelectResponse{
			Response:  err.Error(), // Use the error message as the response
			Protobufs: nil,         // Initialize the protobufs slice
		}, err
	}

	var selectQuery = selectionQuery(columnType, ks, tableName, column, constraint)

	// Execute query
	iter := session.Query(selectQuery).Iter()

	// Prints out total # of rows updated
	if *debug {
		rowCount := iter.NumRows()
		fmt.Println("Selected:", rowCount, "entries")
	}

	// Initialize a new response
	response := &Routes.ProtobufSelectResponse{
		Response:  "",
		Protobufs: nil, // initialize the protobufs slice
	}

	// Iterate through each returned row row
	for {
		// Initialize a new map for each row
		columnValues := make(map[string]interface{})

		// Attempt to scan the values for the current row
		if !iter.MapScan(columnValues) {
			break // Exit the loop if there are no more rows
		}

		if rawBytes, ok := columnValues["serial_msg"].([]byte); ok {
			response.Protobufs = append(response.Protobufs, rawBytes)
		}
	}

	return response, nil
}

/*
* Update handles UPDATE requests based on specified conditions.

* Parameters:
*   - ctx: A context object for the request
*   - request: A protobuf containing information about the keyspace, table, and constraint
*
* Returns:
*   - *ProtobufErrorResponse: A response containing error information if an error occurs.
 */
func (s *server) Update(ctx context.Context, request *Routes.ProtobufUpdateRequest) (*Routes.ProtobufErrorResponse, error) {
	// Lock to ensure only one goroutine is executed at a time
	s.mu.Lock()
	defer s.mu.Unlock()

	// Retrieves data from received protobuf
	tableName := request.Table
	ks := request.Keyspace
	column := request.Column
	constraint := request.Constraint
	newValue := request.NewValue

	// Create a session to interact with the database
	session, err := s.cluster.CreateSession()
	if err != nil {
		log.Printf("Error opening cluster session")
		return nil, err
	}
	defer session.Close()

	// Creates index
	if err := createIndex(session, ks, tableName, column); err != nil {
		return &Routes.ProtobufErrorResponse{Errs: []string{err.Error()}}, err
	}

	// Gets the value type of the column being used
	columnType, err := getColumnType(session, ks, tableName, column)
	if err != nil {
		log.Printf("Error getting column type: %s", err)
		return &Routes.ProtobufErrorResponse{Errs: []string{err.Error()}}, err
	}

	var selectQuery = selectionQuery(columnType, ks, tableName, column, constraint)

	// Execute query and store results in an interator
	iter := session.Query(selectQuery).Iter()

	// Declare a variable to store the ids and keep count
	var idValue string
	counter := 0

	// Loops through returned IDs
	for iter.Scan(&idValue) {
		// Construct UPDATE query with parameter binding (id)
		var updateQuery string

		// Changes query based on type
		switch columnType {
		case "text", "blob", "boolean", "varchar":
			updateQuery = fmt.Sprintf("UPDATE testks.EducationData SET %s = '%s' WHERE serial_msg = ?", column, newValue)
		default:
			updateQuery = fmt.Sprintf("UPDATE testks.EducationData SET %s = %s WHERE serial_msg = ?", column, newValue)
		}

		// Execute UPDATE query
		if updateErr := session.Query(updateQuery, idValue).Exec(); updateErr != nil {
			log.Printf("Error updating data: %s\n query: %s", updateErr, updateQuery)
			return &Routes.ProtobufErrorResponse{Errs: []string{updateErr.Error()}}, updateErr
		}

		// Increment Counter
		counter++
	}

	// Check for errors from the iteration
	if err := iter.Close(); err != nil {
		log.Printf("Error iterating over result: %s\n query: %s", err, selectQuery)
		return &Routes.ProtobufErrorResponse{Errs: []string{err.Error()}}, err
	}

	// Prints out total # of rows updated
	if *debug {
		fmt.Println("Updated", counter, "entries")
	}

	// Returns empty response (No errors encountered)
	return &Routes.ProtobufErrorResponse{}, nil
}

/*
 * Delete handles DELETE requests based on specified conditions.
 *
 * This method deletes from the table based on the provided condition
 *
 * Parameters:
 *   - ctx: A context object for the request
 *   - request: A protobuf containing information about the keyspace and table to drop
 *
 * Returns:
 *   - *ProtobufErrorResponse: A response containing error information if an error occurs.
 */
func (s *server) Delete(ctx context.Context, request *Routes.ProtobufDeleteRequest) (*Routes.ProtobufErrorResponse, error) {
	// Lock to ensure only one goroutine is executed at a time
	s.mu.Lock()
	defer s.mu.Unlock()

	// Retrieves data from received protobuf
	tableName := request.Table
	ks := request.Keyspace
	column := request.Column
	constraint := request.Constraint

	// Create a session to interact with the database
	session, err := s.cluster.CreateSession()
	if err != nil {
		log.Printf("Error opening cluster session")
		return nil, err
	}
	defer session.Close()

	if err := createIndex(session, ks, tableName, column); err != nil {
		return &Routes.ProtobufErrorResponse{Errs: []string{err.Error()}}, err
	}

	// Gets the value type of the column being used
	// Might remove this if we change how the condition is passed in
	columnType, err := getColumnType(session, ks, tableName, column)
	if err != nil {
		log.Printf("Error getting column type: %s", err)
		return &Routes.ProtobufErrorResponse{Errs: []string{err.Error()}}, err
	}

	var selectQuery = selectionQuery(columnType, ks, tableName, column, constraint)

	// Execute query and store results in an interator
	iter := session.Query(selectQuery).Iter()

	// Declare a variable to store the ids and keep count
	var idValue string
	counter := 0

	// Loops through returned IDs
	for iter.Scan(&idValue) {
		// Construct DELETE query with parameter binding (id)
		deleteQuery := "DELETE FROM testks.EducationData WHERE serial_msg = ?"

		// Execute DELETE query
		if deleteErr := session.Query(deleteQuery, idValue).Exec(); deleteErr != nil {
			log.Printf("Error deleting data: %s\n query: %s", deleteErr, deleteQuery)
			return &Routes.ProtobufErrorResponse{Errs: []string{deleteErr.Error()}}, deleteErr
		}

		// Increment Counter
		counter++
	}

	// Check for errors from the iteration
	if err := iter.Close(); err != nil {
		log.Printf("Error iterating over result: %s\n query: %s", err, selectQuery)
		return &Routes.ProtobufErrorResponse{Errs: []string{err.Error()}}, err
	}

	// Prints out total # of rows deleted
	if *debug {
		fmt.Println("Deleted", counter, "entries")
	}

	// Returns empty response (No errors encountered)
	return &Routes.ProtobufErrorResponse{}, nil
}

/*
 * Handle DROP TABLE requests
 *
 * This method drops a table from the  database based on the provided request
 *
 * Parameters:
 *   - ctx: A context object for the request
 *   - request: A protobuf containing information about the keyspace and table to drop
 *
 * Returns:
 *   - *ProtobufErrorResponse: A response containing error information if an error occurs.
 */
func (s *server) DropTable(ctx context.Context, request *Routes.ProtobufDroptableRequest) (*Routes.ProtobufErrorResponse, error) {
	// Lock to ensure only one goroutine is executed at a time
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a session to interact with the database
	session, err := s.cluster.CreateSession()
	if err != nil {
		log.Printf("Error opening cluster session")
		return nil, err
	}
	defer session.Close()

	// Extract keyspace and table from the request
	keyspace := request.Keyspace
	table := request.Table

	// Construct DROP TABLE query with keyspace and table
	dropTableQuery := fmt.Sprintf("DROP TABLE %s.%s;", keyspace, table)

	// Execute DROP TABLE query
	if dropTableErr := session.Query(dropTableQuery).Exec(); dropTableErr != nil {
		log.Printf("Error truncating data: %s\n query: %s", dropTableErr, dropTableQuery)
		return &Routes.ProtobufErrorResponse{Errs: []string{dropTableErr.Error()}}, dropTableErr
	}

	// Prints out total # of rows deleted
	if *debug {
		fmt.Println("Dropped", table, "table")
	}

	return &Routes.ProtobufErrorResponse{}, nil
}

/* Helper Methods */

/* Creates an index as needed for update and delete methods
*
 * Parameters:
 *   - session: A session for executing select query
 *   - keyspace: The keyspace of the target table
 *   - tableName: The name of the target table
 *   - column: The name of the column for the type
 *
 * Returns:
 *	 - error: returns error. No string response needed
*/
func createIndex(session *gocql.Session, keyspace, tableName, column string) error {
	indexQuery := fmt.Sprintf("CREATE INDEX IF NOT EXISTS ON %s.%s (%s)", keyspace, tableName, column)
	if err := session.Query(indexQuery).Exec(); err != nil {
		log.Printf("Error creating index: %s\n query: %s", err, indexQuery)
		return err
	}
	return nil
}

/* Executes selection query for update and delete methods
*
 */
func selectionQuery(columnType string, ks string, tableName string, column string, constraint string) string {
	var returnQuery string

	switch columnType {
	case "text", "blob", "boolean", "varchar":
		// Selects all the id's that meet the condition
		returnQuery = fmt.Sprintf("SELECT serial_msg FROM %s.%s WHERE %s = '%s'", ks, tableName, column, constraint)
	case "int", "bigint", "float", "double", "uuid":
		returnQuery = fmt.Sprintf("SELECT serial_msg FROM %s.%s WHERE %s = %s", ks, tableName, column, constraint)
	}
	return returnQuery
}

/*
 * getColumnType retrieves the type of a specified column in a Cassandra table
 * by querying the system_schema.columns table
 *
 */
func getColumnType(session *gocql.Session, keyspace, tableName, columnName string) (string, error) {
	// Selects all the columns and associated types
	iter := session.Query("SELECT column_name, type FROM system_schema.columns WHERE keyspace_name = ? AND table_name = ?", keyspace, tableName).Iter()

	// Loops through returned column names to look for a matching one
	var fetchedColumnName, columnType string
	for iter.Scan(&fetchedColumnName, &columnType) {
		if fetchedColumnName == columnName {
			return columnType, nil
		}
	}

	// If nothing returned or an error is returned return an error
	if err := iter.Close(); err != nil {
		return "", err
	}

	// If didn't error but nothing found returns error
	return "", fmt.Errorf("Column %s not found in table %s", columnName, tableName)
}

/*
 * convertKindAndValueToCQL maps protobuf.Kind and protobuf.Value to CQL data types.
 *
 * It takes a protobuf field type and value, and returns the corresponding CQL data type
 * as a string and the converted value as an interface{}
 *
 * Parameters:
 *   - fieldType: protoreflect.Kind, the protobuf field type.
 *   - value: protoreflect.Value, the value of the protobuf field.
 *
 * Returns:
 *   - string: The CQL data type.
 *   - interface{}: The converted value.
 *
 * TODO: Handle Unknown or other types.
 */
func convertKindAndValueToCQL(fieldType protoreflect.Kind, value protoreflect.Value) (string, interface{}) {
	switch fieldType {
	case protoreflect.BoolKind:
		return "BOOLEAN", value.Bool()
	case protoreflect.EnumKind, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "INT", value.Int()
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "BIGINT", value.Int()
	case protoreflect.FloatKind:
		return "FLOAT", value.Float()
	case protoreflect.DoubleKind:
		return "DOUBLE", value.Float()
	case protoreflect.StringKind:
		return "TEXT", value.String()
	case protoreflect.BytesKind:
		return "BLOB", value.Bytes()
	default:
		return "UNKNOWN", nil
	}
}

/*
 * Initializes the gRPC server and starts the server to listen for incoming connections.
 *
 * Flags:
 *   - port: Specifies the port on which the gRPC server listens.
 *   - debug: Enables or disables debug mode.
 *
 * Command Line Usage:
 *   - To override the default port (50051), use the -port flag followed by the desired port number.
 *     Example: -port=9090
 *   - To enable debug mode, use the -debug flag.
 *     Example: -debug=true
 *
 * TODO: Create error checking for cmd
 */
func main() {
	// Parses the flags (default if not given)
	flag.Parse()

	// Creates new server
	dbserver := NewServer()

	// Override the port if provided as a command line argument
	if flag.Parsed() {
		// Checks the port flag
		if portFlag := flag.Lookup("port"); portFlag != nil {
			portValue, err := strconv.Atoi(portFlag.Value.String())
			if err != nil {
				log.Fatalf("failed to parse port: %v", err)
			}
			port = &portValue
		}
		// Checks the debug flag
		if debugFlag := flag.Lookup("debug"); debugFlag != nil {
			debugValue, err := strconv.ParseBool(debugFlag.Value.String())
			if err != nil {
				log.Fatalf("failed to parse debug: %v", err)
			}
			debug = &debugValue
		}
	}

	// Create a TCP listener on port var
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	// Register the DataRoute service implementation with the server
	Routes.RegisterDBGenericServer(grpcServer, dbserver)

	// Start the gRPC server as a goroutine
	go func() {
		log.Printf("Server listening at %v", listener.Addr())
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Blocks to keep the server running
	select {}
}
