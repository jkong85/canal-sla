This testbed is similar to the following:
ovs + docker with vxlan [http://www.cnblogs.com/yuuyuu/p/5180827.html]

Start two VMs on host with the bridge mode (copy the interfaces file to create bridge on host)
Then on each VM, run the ovs_test to create the vxlan network 
Different pods between those two VMs commumicate with vxlan

