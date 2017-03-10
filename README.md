Qos for the container network

#Description
We develop the isolation solution for the container network by guaranting and limiting the rate and providing isolated traffic in both of the inbound and outbound direction. It covers several different scenarios, including linux bridge, ovs, ovs-dpdk, macvlan.
(Notes: inbound/outbound is based on the Pod, e.g., the traffic from the Pod to the outside is the outbound, otherwise, the traffic from the outside to the pod is called the inbound traffic)

#Topology
We consider a general topology of the container network based on the linux bridge:
 ---------------------------    ---------------------------
|                      VM0  |  |                      VM1  |
|    Pod0       Pod1        |  |    Pod0       Pod1        |
|        \     /            |  |        \     /            |
|         \   /             |  |         \   /             |
|         Bridge            |  |         Bridge            |
|           |               |  |           |               |
|           |               |  |           |               |
|        VM eth0            |  |        VM eth0            |
 ---------------------------    ---------------------------
            |                               |
            | vxlan                         | vxlan
             -------------------------------

#Main idea
 - Outbound policy: We configure the TBF on each Pod's interafce to limit the outbound rate, as well as configuring the HTB on VM eth0 to guarantee different traffic's bandwidth, priority and providing the bandwidth sharing.

 - Inbound policy: We configure the TBF on the OVS veth interafce connected each Pod to limit the inbound rate, as well as configuring the HTB on the Bridge to guarantee different traffic's bandwidth, priority and providing the bandwidth sharing for the traffics to different Pods.

#Logic
  - Read the Qos policy from ETCD server by while loop
  - Configure the Qos policy on different interface, including the Pod, Bridge and VM

#Solution
    - Problem 1: Cannot configure the TC on VM eth0 to control the Vxlan pkg
        - Solution: configure the bpf directly on the VM's eth0 
        The bpf code is as following:
        ```
        ldh        [12]
        jneq       #0x800,fail
        ldb        [23]
        jneq       #0x11,fail 
        ldh        [36]
        jneq       #0x12b5,fail 
        ld         [76]
        jneq       #0x0a0a1105,fail
        ret        #262144
        fail:   
                ret     #0
        ```
        Then use the bytecode tools bpf_asm to generate the filter. [https://github.com/jkong85/bpftools/tree/master/linux_tools]
#Summary


#Limitation

