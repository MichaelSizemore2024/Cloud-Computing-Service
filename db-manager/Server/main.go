package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	// CQL Imports
	"github.com/gocql/gocql"
	"github.com/joho/godotenv"

	// GRPC Imports
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	Routes "dbmanager/internal/common" // Import the generated code from protofiles
	Log "dbmanager/internal/log"       // Import Logging for the server
	Query "dbmanager/internal/query"   // Import Server Interface
)

// Initializes new server with CRUD interface
func NewServer(node string) *Query.Server {
	return &Query.Server{
		Cluster: initDatabaseCluster(node),
	}
}

// initDatabaseCluster initializes the Cassandra/ScyllaDB cluster configuration.
func initDatabaseCluster(node string) *gocql.ClusterConfig {
	cluster := gocql.NewCluster(node) // Add ScyllaDB node IP or name here
	cluster.Consistency = gocql.One   // Set the consistency level
	return cluster
}

/*
 * Initializes the gRPC server and starts the server to listen for incoming connections.
 */
func main() {
	// Reads in config file
	err := godotenv.Load("internal/config/app.env")
	if err != nil {
		Log.Logger.Info("error loading app.env config file")
		fmt.Println("error loading app.env config file", err)
		return
	}

	// Creates new server instance
	dbserver := NewServer(os.Getenv("NODE"))

	// Initializes Logger
	Log.NewLogger()

	// Gets port from config
	port, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		Log.Logger.Info("error readining in port number")
		fmt.Println("error readining in port number", err)
		return
	}

	// Create a TCP listener on port var
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		Log.Logger.Info("failed to listen on port")
		fmt.Println("failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Register the DataRoute service implementation with the server
	Routes.RegisterDBGenericServer(grpcServer, dbserver)

	// Start the gRPC server as a goroutine
	go func() {
		log.Printf("Server listening at %v", listener.Addr())
		Log.Logger.Info("application started successfully")
		if err := grpcServer.Serve(listener); err != nil {
			Log.Logger.Info("failed to serve")
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Blocks to keep the server running
	select {}
}
