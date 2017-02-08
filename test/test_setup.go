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

	if *flag_clear {
		return
	}
	// create the ovs bridge
	create_ovs_bridge()

	// try to create many dockers
	create_containers()

	// setup the vxlan network (ip start from 10.1.2.3  there are 3 containers)
	create_vxlan_network()
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
}

func clear_vm_config(dev string) {
	fmt.Println("start to clear docker0 config of virtual machine")
	/*
		# delete the existing first
		 sudo ovs-vsctl del-br vxbr
		 sudo tc qdisc del dev docker0 root
	*/
	cmd := "sudo tc qdisc del dev " + dev + " root"
	exe_cmd_full(cmd)
}

func create_vxlan_network() {
	fmt.Println("start to create vxlan network of containrs")
	// get the containder id list
	cmd := "docker"
	args := []string{"ps", "-q"}
	ids := exe_cmd(cmd, args)
	//fmt.Println("ids is ", ids)
	count := start_ip
	// for each containers, do the following
	for _, cid := range strings.Split(ids, "\n") {
		if cid == "" {
			continue
		}

		// IP address +1 , the first is allocated to the bridge
		count = count + 1

		// get the pid
		fmt.Println("cid is: ", cid)
		cmd = "docker"
		args = []string{"inspect", "-f", "{{.State.Pid}}", cid}
		pid := strings.Trim(exe_cmd(cmd, args), "\n")

		// add eth0 for container and bind it to the ovs-bridge
		ovs_docker_cmd := "./ovs-docker add-port vxbr eth0 " + cid
		exe_cmd_full(ovs_docker_cmd)

		// config eth0 ip address of the container by nsenter
		ip := ip_net + strconv.Itoa(count)
		eth_cmd := "nsenter -t " + pid + " -n ifconfig eth0 " + ip + "/24"
		exe_cmd_full(eth_cmd)
		// set mtu for large file trans.
		eth_cmd = "nsenter -t " + pid + " -n ifconfig eth0 mtu 1450"
		exe_cmd_full(eth_cmd)
		// start the ftp server
		//eth_cmd = "nsenter -t " + pid + " -n /usr/sbin/vsftpd & "
		//exe_cmd_full(eth_cmd)

		// change the container name to IPaddress related
		cmd = "docker"
		args = []string{"rename", cid, ip_prefix + "_" + strconv.Itoa(count)}
		exe_cmd(cmd, args)

	}
}

func create_containers() {
	number := number_container
	for number > 0 {
		create_docker_cmd := "docker run --net=none --privileged=true -dit " + image + " bash"
		exe_cmd_full(create_docker_cmd)
		number = number - 1
	}
}

func clear_containers() {
	// clean the containers that existed
	clean_exited_cmd := "docker ps -a | grep 'Exited' | awk '{print $1}' | xargs --no-run-if-empty sudo docker rm -f"
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
	//fmt.Printf("exec cmd out: %s\n", out)
	s := string(out)

	return s
}
