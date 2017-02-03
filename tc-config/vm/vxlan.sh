#NS="sudo nsenter -t $1 -n sudo"
NS="sudo "
$NS tc qdisc add dev $1 root handle 1: htb default 30

$NS tc class add dev $1 parent 1: classid 1:1 htb rate 6mbit burst 15k

#$NS tc class add dev $1 parent 1:1 classid 1:10 htb rate 5mbit burst 15k
$NS tc class add dev $1 parent 1:1 classid 1:10 htb rate 1mbit ceil 6mbit burst 15k
$NS tc class add dev $1 parent 1:1 classid 1:20 htb rate 3mbit ceil 6mbit burst 15k
$NS tc class add dev $1 parent 1:1 classid 1:30 htb rate 2mbit ceil 6mbit burst 15k


$NS tc qdisc add dev $1 parent 1:10 handle 10: sfq perturb 10
$NS tc qdisc add dev $1 parent 1:20 handle 20: sfq perturb 10
$NS tc qdisc add dev $1 parent 1:30 handle 30: sfq perturb 10

#U32="sudo nsenter -t $1 -n sudo tc filter add dev $1 protocol ip parent 1:0 prio 1 u32"
#$U32 match ip protocol 1 0xff flowid 1:10
#$U32 match ip dport 4789 0xffff flowid 1:20

#U32="sudo nsenter -t $1 -n sudo tc filter add dev $1 parent 1:0 prio 1 u32"
#U32="sudo tc filter add dev $1 parent 1:0 prio 1 u32"
#$U32 protocol ip match ip protocol 1 0xff flowid 1:10
#$U32 protocol ip match ip dport 4789 0xffff flowid 1:20

# for all traffic
#sudo tc filter add dev $1 parent 1:0 bpf bytecode "4,40 0 0 12,21 0 1 2048,6 0 0 262144,6 0 0 0" flowid 1:10 
# for all UDP traffic
sudo tc filter add dev $1 parent 1:0 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10

# for Vxlan with dst address : 10.1.2.6
sudo tc filter add dev $1 parent 1:0 bpf bytecode "11,40 0 0 12,21 0 8 2048,48 0 0 23,21 0 6 17,40 0 0 42,69 1 0 2048,6 0 0 0,33 0 0 76,21 0 1 167838213,6 0 0 262144,6 0 0 0," flowid 1:20 


# for dst address : 10.1.2.6
# 11,40 0 0 12,21 0 8 2048,48 0 0 23,21 0 6 17,40 0 0 42,69 1 0 2048,6 0 0 0,32 0 0 80,21 0 1 167838214,6 0 0 262144,6 0 0 0,

#11,40 0 0 12,21 0 8 2048,48 0 0 23,21 0 6 17,40 0 0 42,69 1 0 2048,6 0 0 0,32 0 0 76,21 0 1 167838213,6 0 0 262144,6 0 0 0,
