import grpc
import example_pb2
import example_pb2_grpc

def run():
    # Connect to the gRPC server
    channel = grpc.insecure_channel('localhost:50051')

    # Create a gRPC stub
    stub = example_pb2_grpc.GreeterStub(channel)

    # Prepare a request message
    request = example_pb2.HelloRequest(name='World')

    # Call the gRPC service
    response = stub.SayHello(request)

    # Print the response
    print(f"Received: {response.message}")

if __name__ == '__main__':
    run()
