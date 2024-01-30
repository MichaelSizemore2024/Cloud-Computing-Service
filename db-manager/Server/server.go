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
	"google.golang.org/protobuf/reflect/protoreflect"

	Routes "dbmanager/common" // Import the generated code from protofiles
)

// Variables that can be passed in via command line such as --port=50051 or -port=50051
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

// Handles INSERT requests and returns nothing unless an error is encountered
func (s *server) Insert(ctx context.Context, request *Routes.ProtobufInsertRequest) (*Routes.ProtobufErrorResponse, error) {
	// TODO: maybe use slog package for cleaner/detailed print statements
	// TODO: ((?) might already be done) maybe need to lock during data handling

	messages := []protoreflect.ProtoMessage{} // Create a emtpy list of proto messages
	for _, any := range request.Protobufs {
		msg, err := any.UnmarshalNew() // Unmarshal each Any message to a ProtoMessage
		// Return error to client in a list of strings if encountered
		if err != nil {
			return &Routes.ProtobufErrorResponse{Errs: []string{err.Error()}}, err
		}
		messages = append(messages, msg) // Add new protobuf to list of messages to handle
	}

	// TODO: MAYBE put this in another GOROUTINE (?) - locks might interfere with some of the asynchronous work

	// Loops over all protobufs received and process them
	for _, m := range messages {
		msg_desc := m.ProtoReflect().Descriptor() // find the message descriptor
		fields := msg_desc.Fields()               // get all the fields in the proto message

		queryCols := []string{}                          // list of strings to use strings.Join later for our column names to insert into table
		queryQms := []string{}                           // just to keep track of ?'s in the query to add
		values := []interface{}{gocql.TimeUUID()}        // all the values of each field, input for a variadic function/var, timeUUID for uniqie primary key for now
		queryColTypes := []string{`id UUID PRIMARY KEY`} // for CREATE TABLE, generate the id with the curr datetime to make it unique

		// Loop over each field and infer data
		for i := 0; i < fields.Len(); i++ {
			curField := fields.Get(i)

			queryCols = append(queryCols, string(curField.FullName().Name())) // append the field name to our col list
			queryQms = append(queryQms, "?")                                  // append ?'s for values portion of cql cmds

			convertedType, convertedValue := convertKindAndValueToCQL(curField.Kind(), m.ProtoReflect().Get(curField))
			values = append(values, convertedValue)                                                     // get the value of the field
			queryColTypes = append(queryColTypes, fmt.Sprintf("%s %s", curField.Name(), convertedType)) // NAME TYPE for CQL create table
		}

		s.mu.Lock() // Need to lock here
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
		insertQuery := fmt.Sprint(`INSERT INTO `, ks, `.`, tableName, ` (id, `, strings.Join(queryCols, `, `), `) VALUES (?, `, strings.Join(queryQms, `, `), `)`) // strings.Join is a nice method
		if insertQueryErr := session.Query(insertQuery, values...).Exec(); insertQueryErr != nil {                                                                 // look up Variadic Functions for more on the ... syntax
			log.Printf("Error inserting data: %s\n query: %s", insertQueryErr, insertQuery)
			return &Routes.ProtobufErrorResponse{Errs: []string{insertQueryErr.Error()}}, insertQueryErr
		}
	}

	return &Routes.ProtobufErrorResponse{Errs: []string{}}, nil
}

// Handles UPDATE requests and returns nothing unless an error is encountered
func (s *server) Update(ctx context.Context, request *Routes.ProtobufUpdateRequest) (*Routes.ProtobufErrorResponse, error) {
	s.mu.Lock() // need to lock here
	defer s.mu.Unlock()

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

	// Create index so we can search through contraints (otherwise we can only search via PK)
	indexQuery := fmt.Sprint("CREATE INDEX IF NOT EXISTS ON ", ks, ".", tableName, "(", column, ")")
	if indexErr := session.Query(indexQuery).Exec(); indexErr != nil {
		log.Printf("Error creating index: %s\n query: %s", indexErr, indexQuery)
		return &Routes.ProtobufErrorResponse{Errs: []string{indexErr.Error()}}, indexErr
	}

	// Gets the value type of the column being used
	// Might remove this if we change how the condition is passed in
	columnType, err := getColumnType(session, ks, tableName, column)
	if err != nil {
		log.Printf("Error getting column type: %s", err)
		return &Routes.ProtobufErrorResponse{Errs: []string{err.Error()}}, err
	}

	// Changes the query depending on the type (quotes or no quotes)
	// Will need to add new types as we try different things
	// Might remove this if we change how the condition is passed in, parameter binding doesn't work unless casted
	var selectQuery string
	switch columnType {
	case "text", "blob", "boolean", "varchar":
		// Selects all the id's that meet the condition
		selectQuery = fmt.Sprintf("SELECT id FROM %s.%s WHERE %s = '%s'", ks, tableName, column, constraint)
	case "int", "bigint", "float", "double", "uuid":
		selectQuery = fmt.Sprintf("SELECT id FROM %s.%s WHERE %s = %s", ks, tableName, column, constraint)
	}

	// Execute query
	iter := session.Query(selectQuery).Iter()

	// Declare a variable to store the ids in the loop
	var idValue string

	// Creates counter to keep track # updated
	counter := 0

	// Loops through returned IDs
	for iter.Scan(&idValue) {
		counter++
		// Construct UPDATE query with parameter binding (id)
		var updateQuery string

		switch columnType { //TODO: Protection if using wrong type of data for the column
		case "text", "blob", "boolean", "varchar":
			updateQuery = fmt.Sprintf("UPDATE testks.EducationData SET %s = '%s' WHERE id = %s", column, newValue, idValue)
		default:
			updateQuery = fmt.Sprintf("UPDATE testks.EducationData SET %s = %s WHERE id = %s", column, newValue, idValue)
		}

		// Execute UPDATE query
		if updateErr := session.Query(updateQuery).Exec(); updateErr != nil {
			log.Printf("Error updating data: %s\n query: %s", updateErr, updateQuery)
			return &Routes.ProtobufErrorResponse{Errs: []string{updateErr.Error()}}, updateErr
		}
	}

	// Check for errors from the iteration
	if err := iter.Close(); err != nil {
		log.Printf("Error iterating over result: %s\n query: %s", err, selectQuery)
		return &Routes.ProtobufErrorResponse{Errs: []string{err.Error()}}, err
	}

	// Prints out total # of rows updated, useful for debug can prob be removed in the future
	fmt.Println("Updated", counter, "entries")

	return &Routes.ProtobufErrorResponse{}, nil
}

// Handle DELETE requests
func (s *server) Delete(ctx context.Context, request *Routes.ProtobufDeleteRequest) (*Routes.ProtobufErrorResponse, error) {
	s.mu.Lock() // need to lock here
	defer s.mu.Unlock()

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

	// Create index so we can search through contraints (otherwise we can only search via PK)
	indexQuery := fmt.Sprint("CREATE INDEX IF NOT EXISTS ON ", ks, ".", tableName, "(", column, ")")
	if indexErr := session.Query(indexQuery).Exec(); indexErr != nil {
		log.Printf("Error creating index: %s\n query: %s", indexErr, indexQuery)
		return &Routes.ProtobufErrorResponse{Errs: []string{indexErr.Error()}}, indexErr
	}

	// Gets the value type of the column being used
	// Might remove this if we change how the condition is passed in
	columnType, err := getColumnType(session, ks, tableName, column)
	if err != nil {
		log.Printf("Error getting column type: %s", err)
		return &Routes.ProtobufErrorResponse{Errs: []string{err.Error()}}, err
	}

	// Changes the query depending on the type (quotes or no quotes)
	// Will need to add new types as we try different things
	// Might remove this if we change how the condition is passed in, parameter binding doesn't work unless casted
	var selectQuery string
	switch columnType {
	case "text", "blob", "boolean", "varchar":
		// Selects all the id's that meet the condition
		selectQuery = fmt.Sprintf("SELECT id FROM %s.%s WHERE %s = '%s'", ks, tableName, column, constraint)
	case "int", "bigint", "float", "double", "uuid":
		selectQuery = fmt.Sprintf("SELECT id FROM %s.%s WHERE %s = %s", ks, tableName, column, constraint)
	}

	// Execute query
	iter := session.Query(selectQuery).Iter()

	// Declare a variable to store the ids in the loop
	var idValue string

	// Creates counter to keep track # deleted
	counter := 0

	// Loops through returned IDs
	for iter.Scan(&idValue) {
		counter++
		// Construct DELETE query with parameter binding (id)
		deleteQuery := "DELETE FROM testks.EducationData WHERE id = ?"

		// Execute DELETE query
		if deleteErr := session.Query(deleteQuery, idValue).Exec(); deleteErr != nil {
			log.Printf("Error deleting data: %s\n query: %s", deleteErr, deleteQuery)
			return &Routes.ProtobufErrorResponse{Errs: []string{deleteErr.Error()}}, deleteErr
		}
	}

	// Check for errors from the iteration
	if err := iter.Close(); err != nil {
		log.Printf("Error iterating over result: %s\n query: %s", err, selectQuery)
		return &Routes.ProtobufErrorResponse{Errs: []string{err.Error()}}, err
	}

	// Prints out total # of rows deleted, useful for debug can prob be removed in the future
	if *debug {
		fmt.Println("Deleted", counter, "entries")
	}

	return &Routes.ProtobufErrorResponse{}, nil
}

// Handle DROP TABLE requests
func (s *server) DropTable(ctx context.Context, request *Routes.ProtobufDroptableRequest) (*Routes.ProtobufErrorResponse, error) {
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

	return &Routes.ProtobufErrorResponse{}, nil
}

/* Helper Methods */

// Starts a sesstion and loops through the columns searching for the one matching what is passed in and returns the type of the column in a string the type
func getColumnType(session *gocql.Session, keyspace, tableName, columnName string) (string, error) {
	iter := session.Query("SELECT column_name, type FROM system_schema.columns WHERE keyspace_name = ? AND table_name = ?", keyspace, tableName).Iter()

	var fetchedColumnName, columnType string
	for iter.Scan(&fetchedColumnName, &columnType) {
		if fetchedColumnName == columnName {
			return columnType, nil
		}
	}

	if err := iter.Close(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("Column %s not found in table %s", columnName, tableName)
}

// Helper method that maps protobuf.Kind and protobuf.Value toCQL data types
func convertKindAndValueToCQL(fieldType protoreflect.Kind, value protoreflect.Value) (string, interface{}) {
	switch fieldType {
	//case protoreflect.GroupKind, protoreflect.MessageKind: not sure what do do with these yet
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
	default: // TODO: Handle Unknown or other types
		return "UNKNOWN", nil
	}
}

// TODO: Create error checking for cmd
func main() {
	// Parses the flags (default if not given)
	flag.Parse()

	dbserver := NewServer()

	// Override the port if provided as a command line argument
	if flag.Parsed() {
		if portFlag := flag.Lookup("port"); portFlag != nil {
			portValue, err := strconv.Atoi(portFlag.Value.String())
			if err != nil {
				log.Fatalf("failed to parse port: %v", err)
			}
			port = &portValue
		}
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
