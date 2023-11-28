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
