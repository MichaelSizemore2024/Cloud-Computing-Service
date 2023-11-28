Install VSCode exentions for GO and proto3
Dev Environment Cmds:

Initial Compile Steps:
$ go mod init dbmanager
$ go get google.golang.org/grpc
$ go get google.golang.org/protobuf
$ export PATH="$PATH:$(go env GOPATH)/bin"
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    Routes/dataroute.proto