# Sets the directory the clients are in
CLIENTS_DIR := Clients

# Find all subdirectories in the directory
SUBDIRS := $(wildcard $(CLIENTS_DIR)/*)

# Default targets for 'make'
all: $(SUBDIRS) compile_protos

# Iterate over each directory and call its Makefile if it exists
$(SUBDIRS):
	@if [ -e $@/Makefile ]; then \
        $(MAKE) -C $@; \
    fi

# Add the export command for PATH when using go
export PATH := $(PATH):$(shell go env GOPATH)/bin

# Main targets required for server (compiles everything in proto)
compile_protos:
	@protoc --go_out=internal/common --go-grpc_out=internal/common proto/*.proto

.PHONY: all clean $(SUBDIRS) education_proto email_proto generic_proto weather_proto image_proto

# Default run for the server
run:
	go run Server/main.go

# Clear compiled Protos
clean:
	rm -rf internal/common/*
	$(foreach dir, $(SUBDIRS), $(MAKE) -C $(dir) clean;)
