echo "Current HTB configuration is: "

echo "=======> show qdisc"
sudo tc -s -d qdisc show dev $1 

echo "=======> show class "
sudo tc -s -d class show dev $1

echo "=======> show filter "
sudo tc -s -d filter show dev $1 
