package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	//"reflect"
	"sync"

	"github.com/gocql/gocql" // for scylla

	Routes "dbmanager/common" // Import the generated code

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// variabels that can be passed in via command line such as --port=50051 or -port=50051
var (
	port = flag.Int("port", 50051, "The server port")
)

// Server implements the interface
type server struct {
	Routes.UnimplementedDB_InserterServer
	mu      sync.Mutex           // mutex for thread-safe operations
	cluster *gocql.ClusterConfig // database cluster
}

// Newserver initializes and returns a new server.
func NewServer() *server {
	return &server{
		cluster: initDatabaseCluster(),
	}
}

// initDatabaseCluster initializes the Cassandra/ScyllaDB cluster configuration.
func initDatabaseCluster() *gocql.ClusterConfig {
	cluster := gocql.NewCluster("scylla") // Add ScyllaDB node IP here
	cluster.Consistency = gocql.One       // Set the consistency level
	return cluster
}

// helper method that maps protobuf.Kind and protobuf.Value toCQL data types
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
	//case protoreflect.GroupKind, protoreflect.MessageKind: not sure what do do with these yet
	default: // TODO: unsure what todo with unkown types, CQL lets us create our own types we could do that
		return "UNKNOWN", nil
	}
}

// Handle insert requests
// TODO: maybe use slog package for cleaner/detailed print statements
func (s *server) Insert(ctx context.Context, request *Routes.ProtobufInsertRequest) (*Routes.ProtobufInsertResponse, error) {
	// need to lock during data handling

	messages := []protoreflect.ProtoMessage{} // create emtpy list of proto messages
	for _, any := range request.Protobufs {
		msg, err := any.UnmarshalNew() // Unmarshal each Any message to a ProtoMessage
		if err != nil {
			return &Routes.ProtobufInsertResponse{Errs: []string{err.Error()}}, err // return the error to client in a list of strings
		}
		messages = append(messages, msg) // add new protobuf to list of messages to handle
	}

	// TODO: MAYBE put this in another GOROUTINE ?, locks might interfere with some of the asynchronous work
	// loop over all protobufs received and process them
	for _, m := range messages {
		msg_desc := m.ProtoReflect().Descriptor() // find the message descriptor
		fields := msg_desc.Fields()               // get all the fields in the proto message

		queryCols := []string{}                          // list of strings to use strings.Join later for our column names to insert into table
		queryQms := []string{}                           // just to keep track of ?'s in the query to add
		values := []interface{}{gocql.TimeUUID()}        // all the values of each field, input for a variadic function/var, timeUUID for uniqie primary key for now
		queryColTypes := []string{`id UUID PRIMARY KEY`} // for CREATE TABLE, generate the id with the curr datetime to make it unique

		// loop over each field and infer data
		for i := 0; i < fields.Len(); i++ {
			curField := fields.Get(i)

			queryCols = append(queryCols, string(curField.FullName().Name())) // append the field name to our col list
			queryQms = append(queryQms, "?")                                  // append ?'s for values portion of cql cmds

			convertedType, convertedValue := convertKindAndValueToCQL(curField.Kind(), m.ProtoReflect().Get(curField))
			values = append(values, convertedValue)                                                     // get the value of the field
			queryColTypes = append(queryColTypes, fmt.Sprintf("%s %s", curField.Name(), convertedType)) // NAME TYPE for CQL create table
		}

		s.mu.Lock() // need to lock here
		defer s.mu.Unlock()

		// Create a session to interact with the database
		session, err := s.cluster.CreateSession()
		if err != nil {
			log.Printf("Error opening cluster session")
			return nil, err
		}
		defer session.Close()

		// put together insert query
		tableName := msg_desc.Name()
		ks := request.Keyspace

		// Create the keyspace if it does not exist, note ` vs " matters
		// TODO: CHANGE CLASS AND REPLICATION FACTOR TO ?
		ksQuery := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy','replication_factor': 3}`, ks)
		if ksQueryErr := session.Query(ksQuery).Exec(); ksQueryErr != nil {
			log.Printf("Error creating keyspace: %s\n query: %s", ksQueryErr, ksQuery)
			return &Routes.ProtobufInsertResponse{Errs: []string{ksQueryErr.Error()}}, ksQueryErr
		}

		// now create a table if it does not exist
		// sometimes cql could try to insert a reserved keyword as a column name and it wont work?
		tableQuery := fmt.Sprint(`CREATE TABLE IF NOT EXISTS `, ks, `.`, tableName, ` (`, strings.Join(queryColTypes, `, `), `)`)
		if tableQueryErr := session.Query(tableQuery).Exec(); tableQueryErr != nil {
			log.Printf("Error creating table: %s\n query: %s", tableQueryErr, tableQuery)
			return &Routes.ProtobufInsertResponse{Errs: []string{tableQueryErr.Error()}}, tableQueryErr
		}

		// handle dynamic insert query
		insertQuery := fmt.Sprint(`INSERT INTO `, ks, `.`, tableName, ` (id, `, strings.Join(queryCols, `, `), `) VALUES (?, `, strings.Join(queryQms, `, `), `)`) // strings.Join is a nice method
		if insertQueryErr := session.Query(insertQuery, values...).Exec(); insertQueryErr != nil {                                                                 // look up Variadic Functions for more on the ... syntax
			log.Printf("Error inserting data: %s\n query: %s", insertQueryErr, insertQuery)
			return &Routes.ProtobufInsertResponse{Errs: []string{insertQueryErr.Error()}}, insertQueryErr
		}
	}

	return &Routes.ProtobufInsertResponse{Errs: []string{}}, nil
}

// TODO: Create error checking for cmd
func main() {
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
	Routes.RegisterDB_InserterServer(grpcServer, dbserver)

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
