# Install VSCode exentions for GO and proto3 #

# Initial Compile:
$ go mod init dbmanager
$ export PATH="$PATH:$(go env GOPATH)/bin"

# Cleanup
$ go mod tidy

## Execution ##

# Running the server
$ go run Server/server.go <port>
# i.e for grpc server
$ go run Server/server.go
$ go run Server/server.go -port=50051

# Running Each of the individual clients
$ python education_client.py <csv_file_path> --address <server_address> --port <port_number> <arg>
# i.e for python clients
$ python3 Client/edu_client.py data/education_data.csv
# i.e Add all emails / delete all emails
$ python3 Client/email_client.py data/email_data.csv
$ python3 Client/email_client.py data/email_data.csv --delete



## Handled in Makefile ##

# Compile education.proto
$ python3 -m grpc_tools.protoc --proto_path=./proto --python_out=common --grpc_python_out=common proto/education.proto
$ protoc --go_out=common --go-grpc_out=common proto/education.proto

# Compile email.proto
$ python3 -m grpc_tools.protoc --proto_path=./proto --python_out=common --grpc_python_out=common proto/email.proto
$ protoc --go_out=common --go-grpc_out=common proto/email.proto

# Compile generic.proto
$ python3 -m grpc_tools.protoc --proto_path=./proto --python_out=common --grpc_python_out=common proto/generic.proto
$ protoc --go_out=common --go-grpc_out=common proto/generic.proto



## Scylla EXEC cmds ##

$ cqlsh // 172.20.0.3 9042 // (check logs for specific IP)
$ CREATE KEYSPACE IF NOT EXISTS testks WITH replication = {'class': 'SimpleStrategy','replication_factor': 3};
$ USE testks;
$ CREATE TABLE emaildata (label INT, text TEXT PRIMARY KEY);
$ DROP keyspace testks;