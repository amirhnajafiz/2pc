#!/bin/bash

# running ./main <ip-tables> <mongo-uri> <database> <grpc-port> <test-case-path>
./main ../configs/iptable.txt $(cat ../configs/mongo.txt) global 5001 ../tests/tests.csv
