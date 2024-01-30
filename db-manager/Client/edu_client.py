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

def addAll(csv_file_path, server_address='localhost', server_port=50051):
    # Connect to the gRPC server
    with grpc.insecure_channel(f'{server_address}:{server_port}') as channel:
        # Create a stub (client) for the DBGeneric service
        stub = generic_pb2_grpc.DBGenericStub(channel)  # Updated service name here

        # Read and send data from the CSV file
        for education_data in read_csv(csv_file_path):
            serial_data = education_data.SerializeToString()
            type_url = f"EducationData"
            anypb_msg = any_pb2.Any(value=serial_data, type_url=type_url)

            request = generic_pb2.protobuf_insert_request(
                keyspace="testks",
                protobufs=[anypb_msg]
            )
            response = stub.Insert(request)
            print(f"Server Response: {response.errs}")

def dropTable(server_address='localhost', server_port=50051):
    # Connect to the gRPC server
    with grpc.insecure_channel(f'{server_address}:{server_port}') as channel:
        # Create a stub (client) for the generic service
        stub = generic_pb2_grpc.DBGenericStub(channel)

        # Create a delete request
        droptable_request = generic_pb2.protobuf_droptable_request(
            keyspace="testks",
            table="EducationData"
        )

        # Send the delete request
        response = stub.DropTable(droptable_request)
        print(f"Server Response: {response.errs}")

def delete(server_address='localhost', server_port=50051, table_col=None, col_constraint=None):
    # Connect to the gRPC server
    with grpc.insecure_channel(f'{server_address}:{server_port}') as channel:
        # Create a stub (client) for the generic service
        stub = generic_pb2_grpc.DBGenericStub(channel)

        # Create a delete request
        delete_request = generic_pb2.protobuf_delete_request(
            keyspace="testks",
            table="educationdata",
            column = table_col,
            constraint = col_constraint
        )

        # Send the delete request
        response = stub.Delete(delete_request)
        print(f"Server Response: {response.errs}")
        
def select(server_address='localhost', server_port=50051, table_col=None, col_constraint=None):
    # Connect to the gRPC server
    with grpc.insecure_channel(f'{server_address}:{server_port}') as channel:
        # Create a stub (client) for the generic service
        stub = generic_pb2_grpc.DBGenericStub(channel)

        # Create a delete request
        select_request = generic_pb2.protobuf_select_request(
            keyspace="testks",
            table="educationdata",
            column = table_col,
            constraint = col_constraint
        )

        # Send the delete request
        response = stub.Select(select_request)
        print(f"Server Response: {response.response}")

def update(server_address='localhost', server_port=50051, table_col=None, col_constraint=None, new_value=None):
    # Connect to the gRPC server
    with grpc.insecure_channel(f'{server_address}:{server_port}') as channel:
        # Create a stub (client) for the generic service
        stub = generic_pb2_grpc.DBGenericStub(channel)

        # Create a update request
        update_request = generic_pb2.protobuf_update_request(
            keyspace="testks",
            table="educationdata",
            column = table_col,
            constraint = col_constraint,
            new_value=new_value,
        )

        # Send the update request
        response = stub.Update(update_request)
        print(f"Server Response: {response.errs}")

if __name__ == '__main__':
    # Use argparse to handle command-line arguments
    parser = argparse.ArgumentParser(description='Education gRPC Client')
    parser.add_argument('csv_file_path', help='Path to the CSV file')
    parser.add_argument('--address', default='localhost', help='Address of the gRPC server')  # Add --address argument
    parser.add_argument('--port', type=int, default=50051, help='Port number for the gRPC server')  # Add --port argument

    args = parser.parse_args()

    print("Client listening at port: {}".format(args.port))  # Print the initial message

    while True:
        # Ask the user for a specific flag
        flag = input("Enter a specific flag <AddAll, Delete, DeleteAll, Update, Query, Exit>: ").lower()

        # Check the entered flag and execute the corresponding task
        if flag == 'addall':
            addAll(args.csv_file_path, server_address=args.address, server_port=args.port)
        elif flag == 'delete':
            column = input("Enter a specific column: ")
            constraint = input("Enter a constraint: ")
            delete(server_address=args.address, server_port=args.port, table_col = column, col_constraint = constraint)
        elif flag == 'update':
            column = input("Enter a specific column: ")
            constraint = input("Enter a constraint to change: ")
            updatedVal = input("Enter new value: ")
            update(server_address=args.address, server_port=args.port, table_col = column, col_constraint = constraint, new_value = updatedVal)
        elif flag == 'deleteall':
            dropTable(server_address=args.address, server_port=args.port)
        elif flag == 'query':
            column = input("Enter a specific column: ")
            constraint = input("Enter a constraint: ")
            select(server_address=args.address, server_port=args.port, table_col = column, col_constraint = constraint)
        elif flag == 'exit':
            print("Exited Client.")
            break
        else:
            print("Invalid flag. Please enter a valid flag.")