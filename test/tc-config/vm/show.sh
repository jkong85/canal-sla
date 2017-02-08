echo "Current HTB configuration is: "

echo "=======> show class "
sudo tc -s -d class show dev $1

echo "=======> show filter parent 1:0"
sudo tc filter show dev $1 parent 1:0

echo "=======> show filter parent 1:1"
sudo tc filter show dev $1 parent 1:1

echo "=======> show filter parent 1:2"
sudo tc filter show dev $1 parent 1:2

