# Initial Compile:
    $ go mod init dbmanager
    $ export PATH="$PATH:$(go env GOPATH)/bin"

# Cleanup
    $ go mod tidy

# Execution

# Running the server
$ go run Server/server.go <port> 
    # i.e for grpc server
    $ go run Server/server.go
    $ go run Server/server.go -port=50051
    $ go run Server/server.go -port=50051 -debug=true

# Running the CSV Clients
$ python Clients/csvclients/education_client.py <csv_file_path> --address <server_address> --port <port_number> <arg>
    # i.e for python client
    $ python3 Clients/csvclients/edu_client.py Clients/csvclients/data/education_data.csv

# Running the Weather Client
$ python3 Clients/weatherclient/weather_client.py --address <server_address> --port <port_number> <arg>
    # i.e for python client
    $python3 Clients/weatherclient/weather_client.py 

# Scrapers
    $ python3 Scrapers/xmloutput.py
    $ python3 Scrapers/testscrape.py
    $ python3 Scrapers/statescrape.py
    $ python3 Scrapers/stationscrape.py
    $ python3 Scrapers/weatherscrape.py

## Scylla EXEC cmds ##
 $ cqlsh // 172.20.0.3 9042 // (check logs for specific IP)