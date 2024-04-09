# Initial Compile:
    $ go mod init dbmanager
    $ export PATH="$PATH:$(go env GOPATH)/bin"

# Cleanup
    $ go mod tidy


## Execution ##

# Running the server
$ go run Server/server.go <port> 
    # i.e for grpc server
    $ go run Server/main.go
    $ make run

# Running the server tests
    $ go test server/test.go
    $ make test

# Running the CSV Client
$ python Clients/csvclients/education_client.py <csv_file_path> --address <server_address> --port <port_number> <arg>
    # i.e for python client
    $ python3 Clients/csvclients/edu_client.py Clients/csvclients/data/education_data.csv

# Running the Image client
$ python Clients/imageclient/image_client.py <file_path_for_image_folder> --address <server_address> --port <port_number> <arg> --action <run,deleteall>
    # i.e for python client
    $ python image_client.py data --address 192.168.1.18 --action deleteall

# Running the Weather Client
$ python3 Clients/weatherclient/weather_client.py --address <server_address> --port <port_number> <arg>
    # i.e for python client
    $ python3 Clients/weatherclient/weather_client.py 
    $ python3 Clients/weatherclient/weather_predict.py 

# Scrapers
    $ python3 Clients/weatherclient/Scrapers/xmloutput.py
    $ python3 Clients/weatherclient/Scrapers/statescrape.py
    $ python3 Clients/weatherclient/Scrapers/stationscrape.py
    $ python3 Clients/weatherclient/Scrapers/weatherscrape.py

## Scylla EXEC cmds ##
 $ cqlsh // 172.20.0.3 9042 // (check logs for specific IP)

 # Running the test
 $ cd Server
 $ go test
