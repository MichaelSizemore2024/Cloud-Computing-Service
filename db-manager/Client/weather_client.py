import sys
import os
import time
import grpc
import argparse
from datetime import datetime
import pytz

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

if __name__ == "__main__":
    add_parent_to_path()
    desired_time_zone = pytz.timezone('America/Denver')  # Track Pull Time?
    server_address = 'localhost'
    server_port = 50051

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
                    station_name, temperature, weather = fetch_weather(station)

                    # Check if weather information is available
                    if station_name is not None:
                        # Create a WeatherData Protocol Buffer message
                        weather_data_msg = weather_pb2.WeatherData(
                            state=state,
                            station_name=station_name[:-4], # Removes the State Abbr
                            temp_f=str(temperature), # Must be string, if station is down "N/A" is default
                            weather=weather
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






# Grabs the number of stations for each state
# if __name__ == "__main__":
#     add_parent_to_path()
#     state_names, abbreviations = fetch_state_abbreviations()

#     for state_code in range(len(abbreviations)):
#         print(f"State: {state_names[state_code]} : {abbreviations[state_code]}")

#         # Get stations for Given State
#         station_codes = fetch_station_codes(abbreviations[state_code])
#         print(f"Number of stations in {abbreviations[state_code]}: {len(station_codes)}")


# Individual tests for scraping states weather data
# if __name__ == "__main__":
#     add_parent_to_path()
#     desired_time_zone = pytz.timezone('America/Denver')
#     state_names, abbreviations = fetch_state_abbreviations()

#     # Get state Abbr
#     state_code = int(input("Enter 1-50: ")) - 1  # Zero Index
#     print(f"State: {state_names[state_code]} : {abbreviations[state_code]}")

#     # Enable to fetch number of stations repeatedly
#     # while True:
#     #     station_codes = fetch_station_codes(abbreviations[state_code])
#     #     current_time = datetime.now(desired_time_zone).strftime("%Y-%m-%d %H:%M:%S")
#     #     print(f"{current_time} - Number of station codes: {len(station_codes)}")
#     #     time.sleep(30)  # Sleep for 30 seconds before repeating

#     # Get stations for Given State
#     station_codes = fetch_station_codes(abbreviations[state_code])
#     print(station_codes)

#     # Fetch weather for a specific station code (e.g., KSLC)
#     # Station codes are updated constantly. May show more/less if ran twice
#     weather_result = input("Enter station code: ")
#     station_name, temperature, weather = fetch_weather(weather_result)

    # # Check if weather information is available
    # if station_name is not None:
    #     print(f"Station: {station_name}, Temperature: {temperature} F, Weather: {weather}")
    # else:
    #     print(f"No weather information available for station code {station_name}")