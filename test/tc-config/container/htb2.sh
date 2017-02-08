NS="sudo nsenter -t $1 -n sudo"
$NS tc qdisc add dev eth0 root handle 1: htb default 30

$NS tc class add dev eth0 parent 1: classid 1:1 htb rate 6mbit burst 15k

$NS tc class add dev eth0 parent 1:1 classid 1:10 htb rate 5mbit burst 15k
$NS tc class add dev eth0 parent 1:1 classid 1:20 htb rate 3mbit ceil 6mbit burst 15k
$NS tc class add dev eth0 parent 1:1 classid 1:30 htb rate 1kbit ceil 6mbit burst 15k


$NS tc qdisc add dev eth0 parent 1:10 handle 10: sfq perturb 10
$NS tc qdisc add dev eth0 parent 1:20 handle 20: sfq perturb 10
$NS tc qdisc add dev eth0 parent 1:30 handle 30: sfq perturb 10


U32="sudo nsenter -t $1 -n sudo tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32"
$U32 match ip dport 80 0xffff flowid 1:10
$U32 match ip sport 25 0xffff flowid 1:20
