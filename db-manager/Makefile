.PHONY: all clean

# Add the export command to update PATH
export PATH := $(PATH):$(shell go env GOPATH)/bin

all: education_proto email_proto generic_proto weather_proto

education_proto:
	python3 -m grpc_tools.protoc --proto_path=./proto --python_out=common --grpc_python_out=common proto/education.proto
	protoc --go_out=common --go-grpc_out=common proto/education.proto

email_proto:
	python3 -m grpc_tools.protoc --proto_path=./proto --python_out=common --grpc_python_out=common proto/email.proto
	protoc --go_out=common --go-grpc_out=common proto/email.proto

weather_proto:
	python3 -m grpc_tools.protoc --proto_path=./proto --python_out=common --grpc_python_out=common proto/weather.proto
	protoc --go_out=common --go-grpc_out=common proto/weather.proto

generic_proto:
	python3 -m grpc_tools.protoc --proto_path=./proto --python_out=common --grpc_python_out=common proto/generic.proto
	protoc --go_out=common --go-grpc_out=common proto/generic.proto

clean:
	rm -rf common/*