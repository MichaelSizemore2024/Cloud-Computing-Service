package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/gocql/gocql" // for scylla

	Routes "dbmanager/Routes" // Import the generated code

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb" // import ANY class for future use
)

// variabels that can be passed in via command line such as --port=50051 or -port=50051
var (
	port = flag.Int("port", 50051, "The server port")
)

// Server implements the interface
type server struct {
	Routes.UnimplementedDataRouteServer
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
	cluster.Consistency = gocql.Quorum    // Set the consistency level
	return cluster
}

// createKeyspace creates a keyspace in the database.
func (server *server) createKeyspace(keyspaceName string) error {
	server.mu.Lock()
	defer server.mu.Unlock()

	// Create a session to interact with the database (without specifying keyspace)
	session, err := server.cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Create keyspace query
	keyspaceQuery := fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s
		WITH replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 3
		}
	`, keyspaceName)

	// Execute keyspace query
	if err := session.Query(keyspaceQuery).Exec(); err != nil {
		return err
	}

	fmt.Printf("Keyspace '%s' created successfully\n", keyspaceName)
	return nil
}

// deleteKeyspace deletes a keyspace from the database.
func (server *server) deleteKeyspace(keyspaceName string) error {
	server.mu.Lock()
	defer server.mu.Unlock()

	// Create a session to interact with the database (without specifying keyspace)
	session, err := server.cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Delete keyspace query
	keyspaceQuery := fmt.Sprintf("DROP KEYSPACE IF EXISTS %s", keyspaceName)

	// Execute keyspace query
	if err := session.Query(keyspaceQuery).Exec(); err != nil {
		return err
	}

	fmt.Printf("Keyspace '%s' deleted successfully\n", keyspaceName)
	return nil
}

// GetData implements the GetData RPC method
func (s *server) GetTable(ctx context.Context, req *Routes.TableRequest) (*Routes.TableResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock() // defer runs this once returned from this method

	// Testing Any data type
	any, err := anypb.New(req)
	if err != nil {
		log.Fatalf("Error parsing any datatype: %v", err)
	}
	return &Routes.TableResponse{Name: req.GetName(), Value: any}, nil
}

func main() {
	flag.Parse()

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <create|delete|grpc> <keyspace_name>")
		os.Exit(1)
	}

	action := os.Args[1]
	keyspaceName := os.Args[2]

	app := NewServer()

	switch action {
	case "create":
		err := app.createKeyspace(keyspaceName)
		if err != nil {
			log.Fatalf("Error creating keyspace: %v", err)
		}
	case "delete":
		err := app.deleteKeyspace(keyspaceName)
		if err != nil {
			log.Fatalf("Error deleting keyspace: %v", err)
		}
	case "grpc":
		// Create a TCP listener on port var
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		// Create a new gRPC server
		grpcServer := grpc.NewServer()
		// Register the DataRoute service implementation with the server
		Routes.RegisterDataRouteServer(grpcServer, app)

		// Start the gRPC server as a goroutine
		go func() {
			log.Printf("server listening at %v", listener.Addr())
			if err := grpcServer.Serve(listener); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}()

		// blocks to keep server running
		select {}
	default:
		fmt.Println("Invalid action. Use 'create', 'delete', or 'grpc'.")
		os.Exit(1)
	}
}
