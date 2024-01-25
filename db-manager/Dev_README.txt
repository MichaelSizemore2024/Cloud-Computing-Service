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
    $ go run Server/server.go -port=50051 -debug=true

# Running Each of the individual clients
$ python education_client.py <csv_file_path> --address <server_address> --port <port_number> <arg>
    # i.e for python clients
    $ python3 Client/edu_client.py data/education_data.csv

## Scylla EXEC cmds ##
 $ cqlsh // 172.20.0.3 9042 // (check logs for specific IP)