echo "Current HTB configuration is: "
sudo nsenter -t $1 -n sudo tc -s -d class show dev eth0
