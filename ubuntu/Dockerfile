ARG  CODE_VERSION=latest
FROM ubuntu:${CODE_VERSION}

# Update the package list and upgrade existing packages
RUN apt-get update && \
    apt-get upgrade -y

# Install libraries for compiling
RUN apt-get install -y \
    sudo \
    wget \
    git \
    autoconf libtool pkg-config \
    protobuf-compiler \
    python3 \
    pip

# Create a directory in the root of the container
RUN mkdir -p /root/db-manager

# Change dir to the volume and install GO
WORKDIR /root/db-manager
RUN wget https://go.dev/dl/go1.21.4.linux-amd64.tar.gz

RUN rm -rf /usr/local/go
RUN tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
ENV PATH="$PATH:/usr/local/go/bin"

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Start an interactive shell as default
CMD ["/bin/bash"]   