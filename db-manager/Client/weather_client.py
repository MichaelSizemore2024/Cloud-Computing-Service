import sys
import os
import grpc
import argparse
import pytz
from datetime import datetime

# Change directory to Routes so we can import the protobufs
current_directory = sys.path[0]
routes_directory = current_directory + '/../common'
sys.path.insert(1, routes_directory)

from google.protobuf import any_pb2
import weather_pb2
import generic_pb2
import generic_pb2_grpc

def add_parent_to_path():
    current_dir = os.path.dirname(os.path.realpath(__file__))
    parent_dir = os.path.abspath(os.path.join(current_dir, '..'))
    sys.path.append(parent_dir)

def fetch_state_abbreviations():
    from Scrapers.statescrape import fetch_state_abbreviations
    return fetch_state_abbreviations()

def fetch_station_codes(state_abbr):
    from Scrapers.stationscrape import fetch_station_codes
    return fetch_station_codes(state_abbr)

def fetch_weather(station_code):
    from Scrapers.weatherscrape import fetch_weather_data
    return fetch_weather_data(station_code)

# Input: Fri, 26 Jan 2024 08:53:00 -0700
# Output: 2024-01-26 08:53:00 MST
def convert_to_mst(input_string):
    # Define the format of the input string
    input_format = "%a, %d %b %Y %H:%M:%S %z"

    # Parse the input string to a datetime object
    input_datetime = datetime.strptime(input_string, input_format)

    # Define the MST timezone
    mst_timezone = pytz.timezone('MST')

    # Convert the datetime object to the MST timezone
    mst_datetime = input_datetime.astimezone(mst_timezone)

    # Format the result as a string
    output_string = mst_datetime.strftime("%Y-%m-%d %H:%M:%S %Z")

    return output_string

def run(server_address='localhost', server_port=50051):
    try:
        with grpc.insecure_channel(f'{server_address}:{server_port}') as channel:
            stub = generic_pb2_grpc.DBGenericStub(channel)

            # Get states info
            state_names, abbreviations = fetch_state_abbreviations()
            for state, abbreviation in zip(state_names, abbreviations):
                print(state)
                # Get stations for Given State
                # Station codes are updated constantly. May show more/less if ran twice
                stations = fetch_station_codes(abbreviation)
                for station in stations:
                    # Fetch weather for a specific station code (e.g., KSLC)
                    station_name, temperature, weather, last_update = fetch_weather(station)
                    mtn_datetime = convert_to_mst(last_update)

                    # Check if weather information is available
                    if station_name is not None:
                        # Create a WeatherData Protocol Buffer message
                        weather_data_msg = weather_pb2.WeatherData(
                            state=state,
                            station_name=station_name[:-4], # Removes the State Abbr
                            temp_f=str(temperature), # Must be string, if station is down "N/A" is default
                            weather=weather,
                            last_update=mtn_datetime,
                            station_ID=station
                        )
                        # Serialize the Protocol Buffer message
                        serial_data = weather_data_msg.SerializeToString()

                        # Create an Any message for generic insertion
                        type_url = f"WeatherData"
                        anypb_msg = any_pb2.Any(value=serial_data, type_url=type_url)

                        # Create the protobuf_insert_request
                        request = generic_pb2.protobuf_insert_request(
                            keyspace="weather_data",
                            protobufs=[anypb_msg]
                        )
                        # Send the request to the gRPC server
                        response = stub.Insert(request)
                        # Print server response
                        print(f"Server Response: {response.errs}")
                    else:
                        print(f"No weather information available for station code {station_name}")

    except grpc.RpcError as e:
        print(f"Error communicating with gRPC server: {e}")
        print(f"Code: {e.code()}")
        print(f"Details: {e.details()}")
        print(f"Trailers: {e.trailing_metadata()}")
    except Exception as e:
        import traceback
        traceback.print_exc()
        print(f"An unexpected error occurred: {e}")

if __name__ == "__main__":
    add_parent_to_path()
    # Use argparse to handle command-line arguments
    parser = argparse.ArgumentParser(description='Education gRPC Client')
    parser.add_argument('--address', default='localhost', help='Address of the gRPC server')  # Add --address argument
    parser.add_argument('--port', type=int, default=50051, help='Port number for the gRPC server')  # Add --port argument
    args = parser.parse_args()
    print("Client listening at port: {}".format(args.port))  # Print the initial message
    run(server_address=args.address, server_port=args.port)
