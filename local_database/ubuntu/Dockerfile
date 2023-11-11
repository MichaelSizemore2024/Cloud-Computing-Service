# Use the official Ubuntu 18.04 image as a base (recommended for scylla drivers)
FROM ubuntu:18.04

# Update the package list and upgrade existing packages
RUN apt-get update && \
    apt-get upgrade -y

# Install libraries for compiling
RUN apt-get install -y \
    sudo \
    gcc \
    g++ \
    make \
    cmake \
    wget \
    git

# Download the Scylla C++ Driver
RUN wget https://github.com/scylladb/cpp-driver/releases/download/2.15.2-1/scylla-cpp-driver_2.15.2-1_amd64.deb https://github.com/scylladb/cpp-driver/releases/download/2.15.2-1/scylla-cpp-driver-dev_2.15.2-1_amd64.deb

# Install the Drivers
RUN sudo apt-get install -y ./scylla-cpp-driver_2.15.2-1_amd64.deb ./scylla-cpp-driver-dev_2.15.2-1_amd64.deb

# Create a directory in the root of the container
RUN mkdir -p /root/db-manager

# Copy the Makefile and Connect.cpp into the db-manager directory
COPY Makefile /root/db-manager/Makefile
COPY connect.cpp  /root/db-manager/connect.cpp

# Start an interactive shell as default
CMD ["/bin/bash"]   