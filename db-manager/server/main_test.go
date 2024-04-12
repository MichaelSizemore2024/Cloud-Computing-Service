package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	Routes "dbmanager/internal/common"
	Log "dbmanager/internal/log" // Import Logging for the server
	// Import Logging for the server
)

func TestMain(t *testing.T) {
	// Initializes Logger
	err := os.Chdir(filepath.Join(".."))
	if err != nil {
		panic(err)
	}
	Log.NewLogger()

}

func TestServerInitialization(t *testing.T) {
	// Initialize the server
	server := NewServer("scylla")

	// Check if server is not nil
	assert.NotNil(t, server, "Server should not be nil")

	// Check if cluster is initialized
	assert.NotNil(t, server.Cluster, "Cluster should be initialized")

	// Check if node IP is set correctly
	assert.Equal(t, "scylla", server.Cluster.Hosts[0], "Node IP should be set to localhost")
}

func TestBasicInsert(t *testing.T) {
	// Initialize the server
	server := NewServer("scylla")

	// Create a new context
	ctx := context.TODO()

	// Create a test EmailData proto
	testMessage := &Routes.EmailData{
		Label: 1,
		Text:  "testdata",
	}

	// Marshal the test message to bytes
	serializedMessage, err := proto.Marshal(testMessage)
	assert.NoError(t, err, "Failed to marshal test message")

	// Create an Any message to put in insert proto
	anyMessage := &any.Any{
		TypeUrl: "/EmailData",
		Value:   serializedMessage,
	}

	// Creates a ProtobufInsertRequest
	testRequest := &Routes.ProtobufInsertRequest{
		Keyspace:  "testKeyspace",
		Protobufs: []*any.Any{anyMessage},
	}

	// Call the Insert method and check the response
	response, err := server.Insert(ctx, testRequest)

	// Checks to see if any errors returned
	assert.NoError(t, err, "Unexpected error during Insert")
	assert.NotNil(t, response, "Unexpected nil response")

	// Ensure that the status is "Created"
	assert.Equal(t, Routes.ServerStatus_CREATED, response.Status, "Expected status to be Created")
}

func TestExtremeInsert(t *testing.T) {
	// Initialize the server
	server := NewServer("scylla")

	// Create a new context
	ctx := context.TODO()

	// Generate a large amount of test data (e.g., 100 entries)
	const numEntries = 100
	testData := make([]*Routes.EmailData, numEntries)
	for i := 0; i < numEntries; i++ {
		testData[i] = &Routes.EmailData{
			Label: int32(i),
			Text:  fmt.Sprintf("testdata%d", i),
		}
	}

	// Insert the test data
	for _, data := range testData {
		serializedMessage, err := proto.Marshal(data)
		assert.NoError(t, err, "Failed to marshal test message")

		anyMessage := &any.Any{
			TypeUrl: "/EmailData",
			Value:   serializedMessage,
		}

		testRequest := &Routes.ProtobufInsertRequest{
			Keyspace:  "testKeyspace",
			Protobufs: []*any.Any{anyMessage},
		}

		_, err = server.Insert(ctx, testRequest)
		assert.NoError(t, err, "Unexpected error during Insert")
	}
}

func TestBasicSelect(t *testing.T) {
	// Initialize the server
	server := NewServer("scylla")

	// Create a new context
	ctx := context.TODO()

	selectRequest := &Routes.ProtobufSelectRequest{
		Keyspace:   "testkeyspace",
		Table:      "emaildata",
		Column:     "text",
		Constraint: "testdata",
	}

	// Call the Insert method and check the response
	response, err := server.Select(ctx, selectRequest)

	// Checks to see if any errors returned
	assert.NoError(t, err, "Unexpected error during Insert")
	assert.NotNil(t, response, "Unexpected nil response")

	// Ensure that the status is "Selected"
	assert.Equal(t, Routes.ServerStatus_SELECTED, response.Status, "Expected status to be Selected")
}

func TestBasicUpdate(t *testing.T) {
	// Initialize the server
	server := NewServer("scylla")

	// Create a new context
	ctx := context.TODO()

	updateRequest := &Routes.ProtobufUpdateRequest{
		Keyspace:   "testkeyspace",
		Table:      "emaildata",
		Column:     "label",
		Constraint: "1",
		NewValue:   "1",
	}

	// Call the Insert method and check the response
	response, err := server.Update(ctx, updateRequest)

	// Checks to see if any errors returned
	assert.NoError(t, err, "Unexpected error during Insert")
	assert.NotNil(t, response, "Unexpected nil response")

	// Ensure that the status is "Updated"
	assert.Equal(t, Routes.ServerStatus_UPDATED, response.Status, "Expected status to be Updated")
}

func TestBasicDelete(t *testing.T) {
	// Initialize the server
	server := NewServer("scylla")

	// Create a new context
	ctx := context.TODO()

	deleteRequest := &Routes.ProtobufDeleteRequest{
		Keyspace:   "testkeyspace",
		Table:      "emaildata",
		Column:     "text",
		Constraint: "testdata",
	}

	// Call the Insert method and check the response
	response, err := server.Delete(ctx, deleteRequest)

	// Checks to see if any errors returned
	assert.NoError(t, err, "Unexpected error during Insert")
	assert.NotNil(t, response, "Unexpected nil response")

	// Ensure that the status is "Deleted"
	assert.Equal(t, Routes.ServerStatus_DELETED, response.Status, "Expected status to be Deleted")
}

func TestBasicDropTable(t *testing.T) {
	// Initialize the server
	server := NewServer("scylla")

	// Create a new context
	ctx := context.TODO()

	dropTableRequest := &Routes.ProtobufDroptableRequest{
		Keyspace: "testkeyspace",
		Table:    "emaildata",
	}

	// Call the Insert method and check the response
	response, err := server.DropTable(ctx, dropTableRequest)

	// Checks to see if any errors returned
	assert.NoError(t, err, "Unexpected error during Insert")
	assert.NotNil(t, response, "Unexpected nil response")

	// Ensure that the status is "Deleted"
	assert.Equal(t, Routes.ServerStatus_DELETED, response.Status, "Expected status to be Deleted")
}
