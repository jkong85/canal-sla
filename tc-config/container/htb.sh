sudo nsenter -t $1 -n sudo tc qdisc add dev eth0 root handle 1: htb default 12

sudo nsenter -t $1 -n sudo tc class add dev eth0 parent 1: classid 1:1 htb rate 100kbps ceil 100kbps

sudo nsenter -t $1 -n sudo tc class add dev eth0 parent 1:1 classid 1:10 htb rate 30kbps ceil 100kbps 
sudo nsenter -t $1 -n sudo tc class add dev eth0 parent 1:1 classid 1:11 htb rate 10kbps ceil 100kbps 
sudo nsenter -t $1 -n sudo tc class add dev eth0 parent 1:1 classid 1:12 htb rate 60kbps ceil 100kbps 

#sudo nsenter -t $1 -n sudo tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 \
#    match ip src 172.17.0.2 match ip protocol 1 0xff flowid 1:10

sudo nsenter -t $1 -n sudo tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 \
    match ip src 172.17.0.2 flowid 1:10

sudo nsenter -t $1 -n sudo tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 \
    match ip src 173.17.0.2 match ip dport 80 0xffff flowid 1:11

sudo nsenter -t $1 -n sudo tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 \
    match ip src 1.2.3.4 flowid 1:12


sudo nsenter -t $1 -n sudo tc qdisc add dev eth0 parent 1:10 handle 20: pfifo limit 5
sudo nsenter -t $1 -n sudo tc qdisc add dev eth0 parent 1:11 handle 30: pfifo limit 5
sudo nsenter -t $1 -n sudo tc qdisc add dev eth0 parent 1:12 handle 40: sfq perturb 10

echo "After setup, current configuration is: "
sudo nsenter -t $1 -n sudo tc -s -d class show dev eth0
