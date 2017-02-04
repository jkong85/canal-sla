#!/bin/bash
max=$1
date
for (( i = 1; i < 100; i++ ))
do
    echo "$i"
    ./bwc
    sleep 3 
done
date
