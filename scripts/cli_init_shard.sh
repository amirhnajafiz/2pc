#!/bin/bash

# running ./main shard <path-to-inputs> <path-to-shards> <mongo-db-uri> global
./main shard ../tests/shards.csv $(cat ../configs/local.txt) global
