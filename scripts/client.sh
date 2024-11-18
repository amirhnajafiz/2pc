#!/bin/bash

# running ./main <ip-tables> <mongo-uri> <database>
./main ../configs/iptable.txt $(cat ../configs/mongo.txt) global
