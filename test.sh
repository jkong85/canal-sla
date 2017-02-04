max=$1
date
for (( i=1; i <= $max; ++i ))
do
    echo "$i"
    ./bwctl
    sleep 0.5
done
date
