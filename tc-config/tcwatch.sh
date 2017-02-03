watch -n 1 'sudo tc -s -d qdisc show dev ens3'

echo "=======> show class "
watch -n 1 'sudo tc -s -d class show dev ens3'

echo "=======> show filter parent 1:0"
watch -n 1 'sudo tc -s -d filter show dev ens3 parent 1:0'

echo "=======> show filter parent 1:1"
watch -n 1 'sudo tc -s -d filter show dev ens3 parent 1:1'

echo "=======> show filter parent 1:2"
watch -n 1 'sudo tc -s -d filter show dev ens3 parent 1:2'

