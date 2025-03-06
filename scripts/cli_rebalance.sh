#!/bin/bash

# running ./main rebalance <transactions-count> <mongo-db-uri> global
./main rebalance 1 $(cat ../configs/local.txt) global
