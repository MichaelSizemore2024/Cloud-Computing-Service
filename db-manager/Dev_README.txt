# Install VSCode exentions for GO and proto3 #

# Initial Compile:
$ go mod init dbmanager
$ export PATH="$PATH:$(go env GOPATH)/bin"

# Compile dataroute.proto
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    Routes/dataroute.proto

# Install GO
$ go get google.golang.org/protobuf

# Install gRPC
$ go get google.golang.org/grpc

# Install Scylla driver
$ go get github.com/gocql/gocql

# Cleanup
$ go mod tidy

# Running the server
$ go run Server/server.go <create|delete|grpc> <keyspace_name> <optional flags: -port>
# i.e for grpc server
$ go run Server/server.go grpc -port=50051
# i.e for keyspace
$ go run Server/server.go create TestKeyspaceName

# Running the client
$ go run Client/client.go <optional flags: -name, -addr>

# Scylla Keyspace management
$ go run Server/server.go create ks
$ go run Server/server.go delete ks