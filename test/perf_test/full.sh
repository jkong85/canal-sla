#NS="sudo nsenter -t $1 -n sudo"
echo $(($(date +%s%N)/1000000))

NS="sudo nsenter -t 1 -n "
$NS tc qdisc add dev $1 root handle 1: htb default 30

echo $(($(date +%s%N)/1000000))

$NS tc class add dev $1 parent 1: classid 1:1 htb rate 100mbit burst 15k

echo $(($(date +%s%N)/1000000))

$NS tc class add dev $1 parent 1:1 classid 1:10 htb rate 10mbit ceil 100mbit burst 15k prio 0

echo $(($(date +%s%N)/1000000))

$NS tc class add dev $1 parent 1:1 classid 1:20 htb rate 30mbit ceil 100mbit burst 15k prio 5

echo $(($(date +%s%N)/1000000))

$NS tc class add dev $1 parent 1:1 classid 1:30 htb rate 20mbit ceil 100mbit burst 15k prio 7

echo $(($(date +%s%N)/1000000))


$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10

$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10

$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10
$NS tc filter add dev $1 parent 1:1 protocol ip prio 10 u32 match ip protocol 6 0xff flowid 1:10

echo $(($(date +%s%N)/1000000))
# for all traffic
#sudo tc filter add dev $1 parent 1:0 bpf bytecode "4,40 0 0 12,21 0 1 2048,6 0 0 262144,6 0 0 0" flowid 1:10 
# for all UDP traffic
#sudo tc filter add dev $1 parent 1:0 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10


$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10

$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
$NS tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
echo $(($(date +%s%N)/1000000))

