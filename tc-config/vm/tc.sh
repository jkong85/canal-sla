nic=wlp3s0
tc qdisc add dev $nic root handle 1: htb default 15
tc class add dev $nic parent 1: classid 1:1 htb rate 1000mbit ceil 1000mbit
tc class add dev $nic parent 1:1 classid 1:10 htb rate 80mbit ceil 80mbit prio 0
tc class add dev $nic parent 1:1 classid 1:11 htb rate 80mbit ceil 100mbit prio 1
tc class add dev $nic parent 1:1 classid 1:12 htb rate 20mbit ceil 100mbit prio 2
tc class add dev $nic parent 1:1 classid 1:13 htb rate 20mbit ceil 100mbit prio 2
tc class add dev $nic parent 1:1 classid 1:14 htb rate 10mbit ceil 100mbit prio 3
tc class add dev $nic parent 1:1 classid 1:15 htb rate 30mbit ceil 100mbit prio 3

tc qdisc add dev $nic parent 1:10 handle 100: sfq perturb 10
tc qdisc add dev $nic parent 1:11 handle 110: sfq perturb 10
tc qdisc add dev $nic parent 1:12 handle 120: sfq perturb 10
tc qdisc add dev $nic parent 1:13 handle 130: sfq perturb 10
tc qdisc add dev $nic parent 1:14 handle 140: sfq perturb 10
tc qdisc add dev $nic parent 1:15 handle 150: sfq perturb 10

tc filter add dev $nic parent 1: bpf bytecode "4,40 0 0 12,21 0 1 2048,6 0 0 262144,6 0 0 0" flowid 1:1 #all traffic
tc filter add dev $nic parent 1:1 bpf bytecode "6,40 0 0 12,21 0 3 2048,48 0 0 23,21 0 1 1,6 0 0 262144,6 0 0 0" flowid 1:10 #all icmp traffic
