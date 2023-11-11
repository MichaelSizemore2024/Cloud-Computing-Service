# Use an official C++ base image
FROM gcc:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the contents of your C++ application into the container
COPY . /app

# Build your C++ application
RUN g++ -o your_cpp_executable your_source_file.cpp

# Command to run your C++ application
CMD ["./your_cpp_executable"]