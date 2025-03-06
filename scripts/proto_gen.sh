#!/bin/bash

protoc -I=proto/ --go_out=pkg/ proto/paxos.proto
protoc -I=proto/ --go-grpc_out=pkg/ proto/paxos.proto
protoc -I=proto/ --go_out=pkg/ proto/database.proto
protoc -I=proto/ --go-grpc_out=pkg/ proto/database.proto
