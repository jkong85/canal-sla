tc qdisc add dev ens3 root handle 1: htb default 8001

tc class add dev ens3 parent 1: classid 1:1 htb rate 100mbit ceil 100mbit 
tc class add dev ens3 parent 1:1 classid 1:8001 htb rate 10mbit ceil 100mbit 

tc class add dev ens3 parent 1:1 classid 1:120 htb rate 20mbit ceil 100mbit prio 0 
tc class add dev ens3 parent 1:1 classid 1:119 htb rate 20mbit ceil 100mbit prio 5 
tc class add dev ens3 parent 1:1 classid 1:118 htb rate 20mbit ceil 100mbit prio 7 

tc filter add dev ens3 parent 1: prio 49152 bpf bytecode "10,40 0 0 12,21 0 7 2048,48 0 0 23,21 0 5 17,40 0 0 36,21 0 3 4789,32 0 0 76,21 0 1 167838213,6 0 0 262144,6 0 0 0," flowid 1:120

tc filter add dev ens3 parent 1: prio 49151 bpf bytecode "10,40 0 0 12,21 0 7 2048,48 0 0 23,21 0 5 17,40 0 0 36,21 0 3 4789,32 0 0 76,21 0 1 167838214,6 0 0 262144,6 0 0 0," flowid 1:119


tc filter add dev ens3 parent 1: prio 49150 bpf bytecode "10,40 0 0 12,21 0 7 2048,48 0 0 23,21 0 5 17,40 0 0 36,21 0 3 4789,32 0 0 76,21 0 1 167838215,6 0 0 262144,6 0 0 0," flowid 1:118

