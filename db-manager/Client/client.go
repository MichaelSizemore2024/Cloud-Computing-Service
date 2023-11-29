package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	Routes "dbmanager/Routes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultTableName = "cool default name"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultTableName, "the table to look up in the database")
)

func main() {
	flag.Parse()
	// Set up a connection to the server
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a DataRoute client
	client := Routes.NewDataRouteClient(conn)

	// prevent an err during a faulty context connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call the GetTable RPC
	resp, err := client.GetTable(ctx, &Routes.TableRequest{Name: *name})
	if err != nil {
		log.Fatalf("Error calling GetTable: %v", err)
	}

	// Print the response
	fmt.Printf("Name from server:\n %s\nValue from server(just a test of the any class for now):\n %v\n", resp.GetName(), resp.GetValue())
}
