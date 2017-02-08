echo "Before deleting the HTB configuration"
sudo nsenter -t $1 -n sudo tc -s -d class show dev eth0

sudo nsenter -t $1 -n sudo tc qdisc del dev eth0 root

echo "After deleting the HTB, current configuration is: "
sudo nsenter -t $1 -n sudo tc -s -d class show dev eth0
