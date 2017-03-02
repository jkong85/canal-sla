package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
	ip_net = "10.1.2." // ip address of container
	//start_ip         = 3         // ip_net.start_ip
	//ip_prefix        = "113_"    // container name prefix
	//number_container = 3
)

var number_container int
var start_ip int
var ip_prefix string
var ip_interval int
var ip_vxlan_remote string
var image string

func main() {
	flag_clear := flag.Bool("clear", false, "Only clear the containers")
	flag.IntVar(&number_container, "num", 5, "number of containers to create")
	flag.IntVar(&start_ip, "ip", 0, "the first ip address of containers")
	flag.IntVar(&ip_interval, "interval", 10, "IP address Interval")
	flag.StringVar(&ip_vxlan_remote, "rip", "10.145.240.131", "remote ip for vxlan")
	flag.StringVar(&image, "i", "test", "image for container")
	// first clean the existing containers

	flag.Parse()
	ip_prefix = strconv.Itoa(start_ip)
	start_ip = start_ip*ip_interval + 1

	fmt.Println(" ==> clear the ovs bridge if existed")
	clear_ovs_bridge()
	fmt.Println(" ==> clear the configuration of Docker0 if existed")
	clear_vm_config("docker0")
	fmt.Println(" ==> clear the configuration of Docker0 if existed")
	clear_vm_config("ens3")
	fmt.Println(" ==> clear the containers if existed")
	clear_containers()
	fmt.Println(" ==> clear the macvlan network if existed")
	// after clean the container
	clear_macvlannet()

	if *flag_clear {
		return
	}
	// create the ovs bridge
	create_ovs_bridge()

	// try to create many dockers
	create_containers_macvlannet()
}

func create_ovs_bridge() {
	fmt.Println("start to create ovs bridge")
	/*
		# delete the existing first
		 sudo ovs-vsctl del-br vxbr
		 # add the br
		 sudo ovs-vsctl add-br vxbr
		 sudo ifconfig vxbr 10.1.2.1/24
		 sudo ovs-vsctl add-port vxbr vxlan -- set interface vxlan type=vxlan options:remote_ip=$1
	*/
	cmd := "ovs-vsctl del-br vxbr"
	exe_cmd_full(cmd)

	cmd = " ovs-vsctl add-br vxbr"
	exe_cmd_full(cmd)

	cmd = " ifconfig vxbr " + ip_net + strconv.Itoa(start_ip) + "/24"
	exe_cmd_full(cmd)

	cmd = " ovs-vsctl add-port vxbr vxlan -- set interface vxlan type=vxlan options:remote_ip=" + ip_vxlan_remote
	exe_cmd_full(cmd)

}

func clear_ovs_bridge() {
	fmt.Println("start to clear ovs bridge")
	/*
		# delete the existing first
		 sudo ovs-vsctl del-br vxbr
	*/
	cmd := "ovs-vsctl del-br vxbr"
	exe_cmd_full(cmd)
	// delete the tunnel created by ovs: vxlan-sys-4789
	cmd = "ip link delete vxlan_sys_4789"
	exe_cmd_full(cmd)

}

func clear_vm_config(dev string) {
	fmt.Println("start to clear " + dev + "config of virtual machine")
	/*
		# delete the existing first
		 sudo ovs-vsctl del-br vxbr
		 sudo tc qdisc del dev docker0 root
	*/
	cmd := "sudo tc qdisc del dev " + dev + " root"
	exe_cmd_full(cmd)
}

func create_containers_macvlannet() {
	number := number_container
	for number > 0 {
		//ip := ip_net + strconv.Itoa(number) + "/24"
		// change the container name to IPaddress related : 0_2
		newName := ip_prefix + "_" + strconv.Itoa(number)

		// create macvlan interface
		portName := "int_" + newName
		ovscmd := "ovs-vsctl add-port vxbr " + portName + " -- set interface " + portName + " type=internal"
		exe_cmd_full(ovscmd)
		// create docker macvlan network : docker network create -d macvlan --subnet=10.0.0.0/24 --gateway=10.0.0.1 -o parent=eth0 --ipv6 macvlan0
		subnet := " --subnet=" + ip_net + "0/24 "
		//gateway := " --gateway=" + ip_net + "1 "
		macvlannetName := "macvlannet_" + newName
		dockercmd := "docker network create -d macvlan " + subnet + " -o macvlan_mode=passthru -o parent=" + portName + " " + macvlannetName
		exe_cmd_full(dockercmd)
		// create docker : docker run --net=macvlannet_0_2 --ip=172.16.86.10 -itd jkong85/sharpserver bash
		dockercmd = "docker run " + "--name " + newName + " --net=" + macvlannetName + " -dit " + image + " bash"
		exe_cmd_full(dockercmd)

		number = number - 1
	}
}
func clear_macvlannet() {
	// clean the macvlannet, after stop the container
	//docker network list | grep 'macvlannet_' | awk '{print $2}' | xargs --no-run-if-empty docker network rm
	clean_exited_cmd := " docker network list | grep 'macvlannet_' | awk '{print $2}' | xargs --no-run-if-empty docker network rm"
	exe_cmd_full(clean_exited_cmd)

	/*
		number := number_container
		for number > 0 {
			newName := ip_prefix + "_" + strconv.Itoa(number)
			clean_exited_cmd := "docker network rm macvlannet_" + newName
			exe_cmd_full(clean_exited_cmd)
			number = number - 1
		}
	*/
}

func clear_containers() {

	// clean the containers that existed : docker ps -a | grep 'Exited' | awk '{print $1}' | xargs --no-run-if-empty sudo docker rm -f
	clean_exited_cmd := "docker ps -a | awk '{print $1}' | xargs --no-run-if-empty sudo docker rm -f"
	exe_cmd_full(clean_exited_cmd)

	// get the still up containder id list
	cmd := "docker"
	args := []string{"ps", "-q"}
	ids := exe_cmd(cmd, args)
	fmt.Println("ids is ", ids)
	// for each containers, do the following
	for _, cid := range strings.Split(ids, "\n") {
		if cid == "" {
			continue
		}
		// get the pid
		//fmt.Println("cid is: ", cid)
		cmd = "docker"
		args = []string{"stop", cid}
		exe_cmd(cmd, args)
		// remove the container
		cmd = "docker"
		args = []string{"rm", cid}
		exe_cmd(cmd, args)

		// clear the macvlannet
	}
}

func exe_cmd_full(cmd string) {
	fmt.Println("command is: ", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Println("Error to exec CMD", cmd)
	}
	fmt.Println("cmd output:", string(out))
}

func exe_cmd(cmd string, args []string) string {
	fmt.Println("command is: ", cmd, " ", args)
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Printf("exec cmd error: %s\n", err)
	}
	//log.Println("exec cmd out: %s\n", out)
	s := string(out)

	return s
}
