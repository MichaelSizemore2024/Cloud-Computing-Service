package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	Routes "dbmanager/Routes" // Import the generated code

	"google.golang.org/grpc"
)

// Server implements the interface
type server struct {
	Routes.UnimplementedDataRouteServer
	mu sync.Mutex // mutex for thread-safe operations
}

// GetData implements the GetData RPC method
func (s *server) GetTable(ctx context.Context, req *Routes.TableRequest) (*Routes.TableResponse, error) {
	// Example: return a TableResponse with the received TableRequest
	return &Routes.TableResponse{Name: "Hello!"}, nil
}

func main() {
	// Create a TCP listener on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the DataRoute service implementation with the server
	Routes.RegisterDataRouteServer(grpcServer, &server{})

	// Start the gRPC server
	fmt.Println("gRPC server is listening on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}
