#!/bin/bash

# running ./main shard <path-to-inputs> <path-to-shards> <mongo-db-uri> global
./main shard ../tests/custom/input.csv ../tests/custom/shards.csv $(cat ../configs/mongouri.txt) global
