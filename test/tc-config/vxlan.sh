#NS="sudo nsenter -t $1 -n sudo"
NS="sudo "
$NS tc qdisc add dev $1 root handle 1: htb default 10

$NS tc class add dev $1 parent 1: classid 1:1 htb rate 100mbit burst 15k

#$NS tc class add dev $1 parent 1:1 classid 1:10 htb rate 5mbit burst 15k
$NS tc class add dev $1 parent 1:1 classid 1:10 htb rate 10mbit ceil 60mbit burst 15k prio 0
$NS tc class add dev $1 parent 1:1 classid 1:20 htb rate 30mbit ceil 50mbit burst 15k prio 5
$NS tc class add dev $1 parent 1:1 classid 1:30 htb rate 20mbit ceil 60mbit burst 15k prio 7


#$NS tc qdisc add dev $1 parent 1:10 handle 10: sfq perturb 10
#$NS tc qdisc add dev $1 parent 1:20 handle 20: sfq perturb 10
#$NS tc qdisc add dev $1 parent 1:30 handle 30: sfq perturb 10

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
#sudo tc filter add dev $1 parent 1:0 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10
#sudo tc filter add dev $1 parent 1:1 prio 1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10

# for Vxlan with src address : 10.1.2.3
#sudo tc filter add dev $1 parent 1:0 prio 3 bpf bytecode "11,40 0 0 12,21 0 8 2048,48 0 0 23,21 0 6 17,40 0 0 42,69 1 0 2048,6 0 0 0,32 0 0 76,21 0 1 167838213,6 0 0 262144,6 0 0 0," flowid 1:20 

# for vxlan without VNI, src address 10.1.2.5
sudo tc filter add dev $1 parent 1:0 prio 3 bpf bytecode "10,40 0 0 12,21 0 7 2048,48 0 0 23,21 0 5 17,40 0 0 36,21 0 3 4789,32 0 0 76,21 0 1 167838213,6 0 0 262144,6 0 0 0," flowid 1:20 


# for dst address : 10.1.2.6
# 11,40 0 0 12,21 0 8 2048,48 0 0 23,21 0 6 17,40 0 0 42,69 1 0 2048,6 0 0 0,32 0 0 80,21 0 1 167838214,6 0 0 262144,6 0 0 0,

#11,40 0 0 12,21 0 8 2048,48 0 0 23,21 0 6 17,40 0 0 42,69 1 0 2048,6 0 0 0,32 0 0 76,21 0 1 167838213,6 0 0 262144,6 0 0 0,
