# Install VSCode exentions for GO and proto3 #

# Initial Compile:
$ go mod init dbmanager
$ export PATH="$PATH:$(go env GOPATH)/bin"

# Cleanup
$ go mod tidy

# Compile education.proto
$ python3 -m grpc_tools.protoc --proto_path=./proto --python_out=common --grpc_python_out=common proto/education.proto
$ protoc --go_out=common --go-grpc_out=common proto/education.proto

# Compile email.proto
$ python3 -m grpc_tools.protoc --proto_path=./proto --python_out=common --grpc_python_out=common proto/email.proto
$ protoc --go_out=common --go-grpc_out=common proto/email.proto

# Running the server
$ go run Server/server.go <create|delete|grpc> <arg>
# i.e for grpc server
$ go run Server/server.go grpc
$ go run Server/server.go grpc 50052
# i.e for keyspace
$ go run Server/server.go create TestKeyspaceName
$ go run Server/server.go delete TestKeyspaceName

# Create the table in Scylla (the next block of commands should be run in Scylla's exec)
$ cqlsh
$ USE testkeyspacename;
$ CREATE TABLE email_data (label INT, text TEXT PRIMARY KEY);

# Running Each of the individual clients
$ python education_client.py <csv_file_path> --address <server_address> --port <port_number> <arg>
# i.e for python clients
$ python3 Client/edu_client.py data/education_data.csv
# i.e Add all emails / delete all emails
$ python3 Client/email_client.py data/email_data.csv
$ python3 Client/email_client.py data/email_data.csv --delete


# When complete: you can delete the keyspace that is made
$ DROP keyspace testks;