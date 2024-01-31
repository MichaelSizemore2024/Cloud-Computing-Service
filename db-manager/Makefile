# Sets the directory the clients are in
CLIENTS_DIR := Clients

# Find all subdirectories in the directory
SUBDIRS := $(wildcard $(CLIENTS_DIR)/*)

# Default targets for 'make'
all: $(SUBDIRS) education_proto email_proto generic_proto weather_proto

# Iterate over each directory and call its Makefile if it exists
$(SUBDIRS):
	@if [ -e $@/Makefile ]; then \
        $(MAKE) -C $@; \
    fi

# Add the export command for PATH when using go
export PATH := $(PATH):$(shell go env GOPATH)/bin

# Main targets required for server
education_proto:
	protoc --go_out=common --go-grpc_out=common proto/education.proto

email_proto:
	protoc --go_out=common --go-grpc_out=common proto/email.proto

weather_proto:
	protoc --go_out=common --go-grpc_out=common proto/weather.proto

generic_proto:
	protoc --go_out=common --go-grpc_out=common proto/generic.proto

.PHONY: all clean $(SUBDIRS) education_proto email_proto generic_proto weather_proto

clean:
	rm -rf common/*
	$(foreach dir, $(SUBDIRS), $(MAKE) -C $(dir) clean;)
