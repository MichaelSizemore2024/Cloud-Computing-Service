import sys
import grpc
import csv
import argparse


# Change directory to Routes so we can import the protobufs
current_directory = sys.path[0]
routes_directory = current_directory + '/../common'
sys.path.insert(1, routes_directory)

from google.protobuf.empty_pb2 import Empty  # Update the import statement
import email_pb2_grpc
import email_pb2

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
            email_data = email_pb2.EmailData(
                **{key: get_value(value) for key, value in row.items()}
            )
            yield email_data

def insertAll(csv_file_path, server_address='localhost', server_port=50051):
    # Connect to the gRPC server
    with grpc.insecure_channel(f'{server_address}:{server_port}') as channel:
        # Create a stub (client)
        stub = email_pb2_grpc.EmailServiceStub(channel)

        # Read and send data from the CSV file
        for email_data in read_csv(csv_file_path):
            response = stub.InsertEmail(email_data)
            print(f"Insert: {response.message}")

def deleteAll(server_address='localhost', server_port=50051):
    # Connect to the gRPC server
    with grpc.insecure_channel(f'{server_address}:{server_port}') as channel:
        # Create a stub (client)
        stub = email_pb2_grpc.EmailServiceStub(channel)

        # Create an instance of google.protobuf.Empty
        empty_request = Empty()

        response = stub.DeleteAllEmail(empty_request)
        print(f"Delete: {response.message}")

if __name__ == '__main__':
    # Use argparse to handle command-line arguments
    parser = argparse.ArgumentParser(description='Email gRPC Client')
    parser.add_argument('csv_file_path', help='Path to the CSV file')
    parser.add_argument('--address', default='localhost', help='Address of the gRPC server')
    parser.add_argument('--port', type=int, default=50051, help='Port number for the gRPC server')
    parser.add_argument('--delete', action='store_true', help='Delete emails instead of inserting')

    args = parser.parse_args()

    if args.delete:
        deleteAll(server_address=args.address, server_port=args.port)
    else:
        insertAll(args.csv_file_path, server_address=args.address, server_port=args.port)
