package main

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	Routes "dbmanager/internal/common"
)

func TestBasicInsert(t *testing.T) {
	// Initialize the server
	server := NewServer()

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
}

func TestBasicSelect(t *testing.T) {
	// Initialize the server
	server := NewServer()

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
}

func TestBasicUpdate(t *testing.T) {
	// Initialize the server
	server := NewServer()

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
}

func TestBasicDelete(t *testing.T) {
	// Initialize the server
	server := NewServer()

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
}

func TestBasicDropTable(t *testing.T) {
	// Initialize the server
	server := NewServer()

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
}
