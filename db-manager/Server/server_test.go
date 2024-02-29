package main

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	Routes "dbmanager/common"
)

func TestInsert(t *testing.T) {
	// Initialize the server for testing
	server := NewServer()

	// Create a new context
	ctx := context.TODO()

	// Create a test EmailData proto
	testMessage := &Routes.EmailData{
		Label: 1,
		Text:  "hello jon",
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
