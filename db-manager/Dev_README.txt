# Install VSCode exentions for GO and proto3 #

# Initial Compile:
$ go mod init dbmanager
$ export PATH="$PATH:$(go env GOPATH)/bin"

# Cleanup
$ go mod tidy

# Scraper/Weather Client Install:
$ pip install requests beautifulsoup4
$ pip install lxml
$ pip install pytz #For current time


## Execution ##

# Running the server
$ go run Server/server.go <port> 
    # i.e for grpc server
    $ go run Server/server.go
    $ go run Server/server.go -port=50051
    $ go run Server/server.go -port=50051 -debug=true

# Running Each of the individual clients
# Test CSV Client
$ python education_client.py <csv_file_path> --address <server_address> --port <port_number> <arg>
    # i.e for python clients
    $ python3 Client/edu_client.py data/education_data.csv
# Weather Client
    $ python3 Client/weather_client.py --address <server_address> --port <port_number> <arg>
# Scrapers
    $ python3 Scrapers/xmloutput.py
    $ python3 Scrapers/testscrape.py
    $ python3 Scrapers/statescrape.py
    $ python3 Scrapers/stationscrape.py
    $ python3 Scrapers/weatherscrape.py

## Scylla EXEC cmds ##
 $ cqlsh // 172.20.0.3 9042 // (check logs for specific IP)
