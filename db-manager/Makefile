# Makefile for the connect.cpp program

CC = g++
CFLAGS = -Wall
LIBS = -lscylla-cpp-driver

all: connect

connect: connect.o
	$(CC) $(CFLAGS) -o connect connect.o $(LIBS)

connect.o: connect.cpp
	$(CC) $(CFLAGS) -c connect.cpp

clean:
	rm -f connect *.o

run: connect
	./connect