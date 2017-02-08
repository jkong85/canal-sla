echo "Before deleting the HTB configuration"
sudo tc -s -d class show dev $1 

sudo tc qdisc del dev $1 root

echo "After deleting the HTB, current configuration is: "
sudo tc -s -d class show dev $1 
