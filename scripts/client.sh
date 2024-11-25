#!/bin/bash

# running ./main <ip-tables> <mongo-uri> <database> <grpc-port> <test-case-path>
./main ../configs/hosts.ini $(cat ../configs/local.txt) global 5001 ../tests/tests.csv
