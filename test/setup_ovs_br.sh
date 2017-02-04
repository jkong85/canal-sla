#!/bin/sh

# remove the br at first
sudo ovs-vsctl del-br vxbr
# add the br
sudo ovs-vsctl add-br vxbr

sudo ifconfig vxbr 10.1.2.1/24

sudo ovs-vsctl add-port vxbr vxlan -- set interface vxlan type=vxlan options:remote_ip=$1



