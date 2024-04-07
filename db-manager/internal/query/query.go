package query

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	Routes "dbmanager/internal/common" // Import the generated code from protofiles
	helper "dbmanager/internal/helper" // Import the helper methods
	Log "dbmanager/internal/log"       // Import Logging for the server

	"github.com/gocql/gocql"
	"google.golang.org/protobuf/proto"
)

// Server implements the interface
type Server struct {
	Routes.UnimplementedDBGenericServer
	Mu      sync.Mutex
	Cluster *gocql.ClusterConfig
}

/*
* Handles INSERT requests and stores the provided protobuf messages in the database.
 *
 * Parameters:
 *   - ctx: A context object for the request
 *   - request: A protobuf containing information about the keyspace and table to drop
 *
 * Returns:
 *   - *ProtobufServerResponse: A response containing error information if an error occurs.
*/
func (s *Server) Insert(ctx context.Context, request *Routes.ProtobufInsertRequest) (*Routes.ProtobufServerResponse, error) {
	// Initialize debug counter
	counter := 0

	// Loops over all protobufs received and process them
	for _, any := range request.Protobufs {

		// Unmarshal each Any message to a ProtoMessage
		m, err := any.UnmarshalNew()

		// Return error to client in a list of strings if encountered
		if err != nil {
			Log.Logger.Error("error unmarshaling given protos")
			return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
		}

		msg_desc := m.ProtoReflect().Descriptor() // find the message descriptor
		fields := msg_desc.Fields()               // get all the fields in the proto message

		queryCols := []string{} // list of strings to use strings.Join later for our column names to insert into table
		queryQms := []string{}  // just to keep track of ?'s in the query to add
		serializedAny, err := proto.Marshal(m)

		// Return error if unable to create value used for key
		if err != nil {
			Log.Logger.Error("error marshalling the message")
			return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
		}

		values := []interface{}{serializedAny}                   // all the values of each field, input for a variadic function/var, timeUUID for uniqie primary key for now
		queryColTypes := []string{`serial_msg BLOB PRIMARY KEY`} // for CREATE TABLE, generate the id with the curr datetime to make it unique

		// Adds the serialized message to the table (used a primary key)
		queryCols = append(queryCols, "serial_msg")
		queryQms = append(queryQms, "?")

		// Loop over each field and infer data
		for i := 0; i < fields.Len(); i++ {
			curField := fields.Get(i)

			queryCols = append(queryCols, string(curField.FullName().Name())) // append the field name to our col list
			queryQms = append(queryQms, "?")                                  // append ?'s for values portion of cql cmds

			convertedType, convertedValue := helper.ConvertKindAndValueToCQL(curField.Kind(), m.ProtoReflect().Get(curField))
			values = append(values, convertedValue)                                                     // get the value of the field
			queryColTypes = append(queryColTypes, fmt.Sprintf("%s %s", curField.Name(), convertedType)) // NAME TYPE for CQL create table
		}

		// Lock to ensure only one goroutine is executed at a time
		s.Mu.Lock()
		defer s.Mu.Unlock()

		// Create a session to interact with the database
		session, err := s.Cluster.CreateSession()
		if err != nil {
			Log.Logger.Error("error opening cluster seesion")
			return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
		}
		defer session.Close()

		// Put together insert query
		tableName := msg_desc.Name()
		ks := request.Keyspace

		// Create the keyspace if it does not exist
		ksQuery := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy','replication_factor': 3}`, ks)
		if ksQueryErr := session.Query(ksQuery).Exec(); ksQueryErr != nil {
			Log.Logger.Error("error creating keyspace")
			return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
		}

		// Create a table if it does not exist
		tableQuery := fmt.Sprint(`CREATE TABLE IF NOT EXISTS `, ks, `.`, tableName, ` (`, strings.Join(queryColTypes, `, `), `)`)
		if tableQueryErr := session.Query(tableQuery).Exec(); tableQueryErr != nil {
			Log.Logger.Error("error creating table")
			return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
		}

		// handle dynamic insert query
		insertQuery := fmt.Sprint(`INSERT INTO `, ks, `.`, tableName, ` ( `, strings.Join(queryCols, `, `), `) VALUES (`, strings.Join(queryQms, `, `), `)`)
		if insertQueryErr := session.Query(insertQuery, values...).Exec(); insertQueryErr != nil {
			Log.Logger.Error("error inserting data")
			return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
		}

		// Increment counter
		counter++
	}
	Log.Logger.Info("Inserted " + strconv.Itoa(counter) + " entries")
	return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_CREATED, Errs: []string{}}, nil
}

/*
* Select handles SELECT requests based on specified conditions.

* Parameters:
*   - ctx: A context object for the request
*   - request: A protobuf containing information about the keyspace, table, and constraint
*
* Returns:
*   - *ProtobufSelectResponse: A response containing error information if an error occurs.
 */
func (s *Server) Select(ctx context.Context, request *Routes.ProtobufSelectRequest) (*Routes.ProtobufSelectResponse, error) {
	// Lock to ensure only one goroutine is executed at a time
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Retrieves data from received protobuf
	tableName := request.Table
	ks := request.Keyspace
	column := request.Column
	constraint := request.Constraint

	// Create a session to interact with the database
	session, err := s.Cluster.CreateSession()
	if err != nil {
		Log.Logger.Error("error opening cluster seesion")
		return &Routes.ProtobufSelectResponse{Status: Routes.ServerStatus_FAILED, Errs: err.Error()}, err
	}
	defer session.Close()

	// Create index so we can search through constraints (otherwise we can only search via PK)
	if err := helper.CreateIndex(session, ks, tableName, column); err != nil {
		Log.Logger.Error("error creating index")
		return &Routes.ProtobufSelectResponse{
			Status:    Routes.ServerStatus_FAILED,
			Errs:      err.Error(), // Use the error message as the response
			Protobufs: nil,         // Initialize the protobufs slice
		}, err
	}

	// Gets the value type of the column being used
	columnType, err := helper.GetColumnType(session, ks, tableName, column)
	if err != nil {
		Log.Logger.Error("error getting column types")
		return &Routes.ProtobufSelectResponse{
			Status:    Routes.ServerStatus_FAILED,
			Errs:      err.Error(), // Use the error message as the response
			Protobufs: nil,         // Initialize the protobufs slice
		}, err
	}

	// Create query
	var selectQuery = helper.SelectionQuery(columnType, ks, tableName, column, constraint)

	// Execute query
	iter := session.Query(selectQuery).Iter()

	// Prints out total # of rows updated
	rowCount := iter.NumRows()
	Log.Logger.Info("Selected " + strconv.Itoa(rowCount) + " entries")

	// Initialize a new response
	response := &Routes.ProtobufSelectResponse{
		Status:    Routes.ServerStatus_SELECTED,
		Errs:      "",
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
func (s *Server) Update(ctx context.Context, request *Routes.ProtobufUpdateRequest) (*Routes.ProtobufServerResponse, error) {
	// Lock to ensure only one goroutine is executed at a time
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Retrieves data from received protobuf
	tableName := request.Table
	ks := request.Keyspace
	column := request.Column
	constraint := request.Constraint
	newValue := request.NewValue

	// Create a session to interact with the database
	session, err := s.Cluster.CreateSession()
	if err != nil {
		Log.Logger.Error("error creating cluster session")
		return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
	}
	defer session.Close()

	// Creates index
	if err := helper.CreateIndex(session, ks, tableName, column); err != nil {
		Log.Logger.Error("error creating index")
		return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
	}

	// Gets the value type of the column being used
	columnType, err := helper.GetColumnType(session, ks, tableName, column)
	if err != nil {
		Log.Logger.Error("error getting column type")
		return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
	}

	// Creates selection Query
	var selectQuery = helper.SelectionQuery(columnType, ks, tableName, column, constraint)

	// Execute query and store results in an interator
	iter := session.Query(selectQuery).Iter()

	// Declare a variable to store the ids and keep count
	var idValue string
	counter := 0

	// Loops through returned IDs
	for iter.Scan(&idValue) {
		// Construct UPDATE query with parameter binding (id)
		var updateQuery string

		// Construct the fully qualified table name
		qualifiedTableName := fmt.Sprintf("%s.%s", ks, tableName)

		// Changes query based on type
		switch columnType {
		case "TEXT", "BLOB", "BOOLEAN", "VARCHAR":
			updateQuery = fmt.Sprintf("UPDATE %s SET %s = '%s' WHERE serial_msg = ?", qualifiedTableName, column, newValue)
		default:
			updateQuery = fmt.Sprintf("UPDATE %s SET %s = %s WHERE serial_msg = ?", qualifiedTableName, column, newValue)
		}

		// Execute UPDATE query
		if updateErr := session.Query(updateQuery, idValue).Exec(); updateErr != nil {
			Log.Logger.Error("error updating the table")
			return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
		}

		// Increment Counter
		counter++
	}

	// Check for errors from the iteration
	if err := iter.Close(); err != nil {
		Log.Logger.Error("error encountered while updating")
		return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
	}

	// Prints out total # of rows updated
	Log.Logger.Info("Updated " + strconv.Itoa(counter) + " entries")

	// Returns empty response (No errors encountered)
	return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_UPDATED}, nil
}

/*
 * Delete handles DELETE requests based on specified conditions.
 *
 * Parameters:
 *   - ctx: A context object for the request
 *   - request: A protobuf containing information about the keyspace and table to drop
 *
 * Returns:
 *   - *ProtobufErrorResponse: A response containing error information if an error occurs.
 */
func (s *Server) Delete(ctx context.Context, request *Routes.ProtobufDeleteRequest) (*Routes.ProtobufServerResponse, error) {
	// Lock to ensure only one goroutine is executed at a time
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Retrieves data from received protobuf
	tableName := request.Table
	ks := request.Keyspace
	column := request.Column
	constraint := request.Constraint

	// Create a session to interact with the database
	session, err := s.Cluster.CreateSession()
	if err != nil {
		Log.Logger.Error("error creating cluster session")
		return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
	}
	defer session.Close()

	if err := helper.CreateIndex(session, ks, tableName, column); err != nil {
		Log.Logger.Error("error creating index")
		return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
	}

	columnType, err := helper.GetColumnType(session, ks, tableName, column)
	if err != nil {
		Log.Logger.Error("error getting column type")
		return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
	}

	// Create selection query
	var selectQuery = helper.SelectionQuery(columnType, ks, tableName, column, constraint)

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
			Log.Logger.Error("error deleteing data from table")
			return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
		}

		// Increment Counter
		counter++
	}

	// Check for errors from the iteration
	if err := iter.Close(); err != nil {
		Log.Logger.Error("error encountered while deleteing entries")
		return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
	}

	// Prints out total # of rows updated
	Log.Logger.Info("Deleted " + strconv.Itoa(counter) + " entries")

	// Returns empty response (No errors encountered)
	return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_DELETED}, nil
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
func (s *Server) DropTable(ctx context.Context, request *Routes.ProtobufDroptableRequest) (*Routes.ProtobufServerResponse, error) {
	// Lock to ensure only one goroutine is executed at a time
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Create a session to interact with the database
	session, err := s.Cluster.CreateSession()
	if err != nil {
		Log.Logger.Error("error connecting to culster session")
		return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{err.Error()}}, err
	}
	defer session.Close()

	// Extract keyspace and table from the request
	keyspace := request.Keyspace
	table := request.Table

	// Construct DROP TABLE query with keyspace and table
	dropTableQuery := fmt.Sprintf("DROP TABLE %s.%s;", keyspace, table)

	// Execute DROP TABLE query
	if dropTableErr := session.Query(dropTableQuery).Exec(); dropTableErr != nil {
		Log.Logger.Error("error encountered when dropping table")
		return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_FAILED, Errs: []string{dropTableErr.Error()}}, dropTableErr
	}

	// Prints out total # of rows deleted
	Log.Logger.Info("Dropped " + table + " table")

	return &Routes.ProtobufServerResponse{Status: Routes.ServerStatus_DELETED}, nil
}
