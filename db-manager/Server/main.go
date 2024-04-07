package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/gocql/gocql"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	Routes "dbmanager/internal/common" // Import the generated code from protofiles
	Log "dbmanager/internal/log"
	Query "dbmanager/internal/query" // Import Server Interface
)

// Variables that can be passed in via command line such as -port=50051
var (
	port = flag.Int("port", 50051, "The server port")
)

// Initializes new server with CRUD interface
func NewServer() *Query.Server {
	return &Query.Server{
		Cluster: initDatabaseCluster(),
	}
}

// initDatabaseCluster initializes the Cassandra/ScyllaDB cluster configuration.
func initDatabaseCluster() *gocql.ClusterConfig {
	cluster := gocql.NewCluster("scylla") // Add ScyllaDB node IP or name here (scylla in this instance)
	cluster.Consistency = gocql.One       // Set the consistency level
	return cluster
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
 *
 * TODO: Create error checking for cmd
 */
func main() {
	// Parses the flags (default if not given)
	flag.Parse()

	// Creates new server
	dbserver := NewServer()

	// Initializes Logger
	Log.NewLogger()

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
		Log.Logger.Info("Application Started successfully on port 80")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Blocks to keep the server running
	select {}
}
