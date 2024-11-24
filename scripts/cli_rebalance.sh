#!/bin/bash

# running ./main rebalance <transactions-count> <apply/dry-run> <mongo-db-uri> global
./main rebalance 1 dry-run $(cat ../configs/local.txt) global
