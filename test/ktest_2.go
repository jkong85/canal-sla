package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var number_container int
var start_ip int
var ip_prefix string
var ip_interval int
var ip_net string
var image string

func main() {
	flag_clear := flag.Bool("clear", false, "Only clear the containers")
	flag.IntVar(&number_container, "num", 2, "number of containers to create")
	flag.IntVar(&start_ip, "ip", 100, "the first ip address of containers")
	flag.StringVar(&ip_net, "net", "172.16.50.", "ip net")
	// first clean the existing containers

	flag.Parse()
	ip_prefix = strconv.Itoa(start_ip)

	fmt.Println(" ==> clear the configuration of Docker0 if existed")
	clear_vm_config("docker0")
	fmt.Println(" ==> clear the configuration of Docker0 if existed")
	clear_vm_config("ens3")
	fmt.Println(" ==> clear the containers if existed")
	clear_containers()

	if *flag_clear {
		return
	}

	// try to create many dockers
	create_containers()

	// setup the vxlan network (ip start from 10.1.2.3  there are 3 containers)
	set_container()
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

func set_container() {
	fmt.Println("start to set container")
	// get the containder id list
	cmd := "docker"
	args := []string{"ps", "-q", "-n", strconv.Itoa(number_container)}
	show_cmd := "docker ps -a | grep '80/tcp' | awk '{print $1}' | xargs --no-run-if-empty "
	ids := exe_cmd_full_2(show_cmd)

	s1 := ids
	if last := len(s1) - 1; last >= 0 && s1[last] == '\n' {
		s1 = s1[:last]
	}
	fmt.Println("new ids is ", s1)
	start_ip = 100
	count := start_ip
	// for each containers, do the following
	for _, cid := range strings.Split(s1, " ") {
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

		// delete its orginal ip address
		cmd = "nsenter"
		args = []string{"-t", pid, "-n", "ifconfig", "eth0"}
		ifconfig_output := exe_cmd(cmd, args)

		var eth_cmd string

		for _, line := range strings.Split(ifconfig_output, "\n") {
			if strings.Contains(line, "inet addr") {
				container_ip := strings.Split(strings.Split(line, ":")[1], " ")[0]
				//fmt.Printf("container id: %s; pid: %s; ip addr: %s\n", container_id, container_pid, container_ip)
				eth_cmd = "nsenter -t " + pid + " -n ip addr del " + container_ip + "/24 dev eth0"
				exe_cmd_full(eth_cmd)
			}
		}

		// config eth0 ip address of the container by nsenter
		ip := ip_net + strconv.Itoa(count)
		eth_cmd = "nsenter -t " + pid + " -n ip addr add " + ip + "/24 dev eth0"
		exe_cmd_full(eth_cmd)

		// config eth0 ip address of the container by nsenter
		ip_default := ip_net + "1"
		eth_cmd = "nsenter -t " + pid + " -n ip route add default via " + ip_default
		exe_cmd_full(eth_cmd)

		// set mtu for large file trans.
		//eth_cmd = "nsenter -t " + pid + " -n ifconfig eth0 mtu 1450"
		//exe_cmd_full(eth_cmd)

		// change the container name to IPaddress related
		cmd = "docker"
		name := ip_prefix + "_" + strconv.Itoa(count)
		args = []string{"rename", cid, name}
		exe_cmd(cmd, args)
		// start the httpd
		eth_cmd = "docker exec -d " + name + " /usr/local/apache2/bin/httpd"
		exe_cmd_full(eth_cmd)
	}
}

func create_containers() {
	number := number_container
	for number > 0 {
		create_docker_cmd := "docker run -dit httpd:ssh bash "
		exe_cmd_full(create_docker_cmd)
		number = number - 1
	}
}

func clear_containers() {

	// get the still up containder id list
	cmd := "docker"
	args := []string{"ps", "-q"}
	//ids := exe_cmd(cmd, args)

	// clean the containers that existed
	clean_exited_cmd := "docker ps -a | grep ' httpd ' | awk '{print $1}' | xargs --no-run-if-empty "
	ids := exe_cmd_full_2(clean_exited_cmd)

	s1 := ids
	if last := len(s1) - 1; last >= 0 && s1[last] == '\n' {
		s1 = s1[:last]
	}

	fmt.Println("ids is ", ids)
	fmt.Println("and new ids is ", s1)
	// for each containers, do the following
	for _, cid := range strings.Split(s1, " ") {
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

func exe_cmd_full_2(cmd string) string {
	fmt.Println("command is: ", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Println("Error to exec CMD", cmd)
	}
	fmt.Println("cmd output:", string(out))

	s := string(out)

	return s
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
