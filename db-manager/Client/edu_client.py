import sys
import grpc
import csv
import argparse

# Change directory to Routes so we can import the protobufs
current_directory = sys.path[0]
routes_directory = current_directory + '/../common'
sys.path.insert(1, routes_directory)

from google.protobuf import any_pb2
import education_pb2
import education_pb2_grpc
import generic_pb2
import generic_pb2_grpc

def get_value(value):
    try:
        return int(value)
    except ValueError:
        try:
            return float(value)
        except ValueError:
            return value

def read_csv(file_path):
    with open(file_path, newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            education_data = education_pb2.EducationData(
                **{key: get_value(value) for key, value in row.items()}
            )
            yield education_data

def run(csv_file_path, server_address='localhost', server_port=50051):
    # Connect to the gRPC server
    with grpc.insecure_channel(f'{server_address}:{server_port}') as channel:
        # Create a stub (client)
        stub = generic_pb2_grpc.DB_InserterStub(channel)

        # Read and send data from the CSV file
        for education_data in read_csv(csv_file_path):
            serial_data = education_data.SerializeToString()
            type_url = f"EducationData" # this is needed if we define a package in the proto then it  would be the packagename.EmailData or packagename.email_pb2.EmailData
            anypb_msg = any_pb2.Any(value=serial_data, type_url=type_url)

            request = generic_pb2.protobuf_insert_request(
                keyspace="testks",
                protobufs=[anypb_msg] 
            )
            response = stub.Insert(request)
            print(f"Server Response: {response.errs}")

if __name__ == '__main__':
    # Use argparse to handle command-line arguments
    parser = argparse.ArgumentParser(description='Education gRPC Client')
    parser.add_argument('csv_file_path', help='Path to the CSV file')
    parser.add_argument('--address', default='localhost', help='Address of the gRPC server')  # Add --address argument
    parser.add_argument('--port', type=int, default=50051, help='Port number for the gRPC server')  # Add --port argument

    args = parser.parse_args()

    # Runs the program with the provided arguments
    run(args.csv_file_path, server_address=args.address, server_port=args.port)