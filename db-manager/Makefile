.PHONY: all clean

# Add the export command to update PATH
export PATH := $(PATH):$(shell go env GOPATH)/bin

all: education_proto email_proto generic_proto weather_proto

education_proto:
	protoc --go_out=common --go-grpc_out=common proto/education.proto

email_proto:
	protoc --go_out=common --go-grpc_out=common proto/email.proto

weather_proto:
	protoc --go_out=common --go-grpc_out=common proto/weather.proto

generic_proto:
	protoc --go_out=common --go-grpc_out=common proto/generic.proto

clean:
	rm -rf common/*