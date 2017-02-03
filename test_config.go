package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	create_dockers(10)
}

func create_dockers(number int) {
	// try to create many dockers

	// get the containder id list
	cmd := "docker"
	args := []string{"ps", "-q"}
	ids := exe_cmd(cmd, args)
	count := 3
	// for each containers, do the following
	for _, cid := range strings.Split(ids, "\n") {
		if cid == "" {
			continue
		}
		// get the pid
		cmd = "docker"
		args = []string{"inspect", "-f", "{{.State.Pid}}", cid}
		pid := strings.Trim(exe_cmd(cmd, args), "\n")
		// add eth0 for container and bind it to the ovs-bridge
		ovs_docker_cmd := "./ovs-docker add-port vxbr eth0 " + cid
		exe_cmd_full(ovs_docker_cmd)
		// config eth0 ip address of the container by nsenter
		ip := "10.1.2." + strconv.Itoa(count)
		eth_cmd := "nsenter -t " + pid + " -n ifconfig eth0 " + ip + "/24"
		exe_cmd_full(eth_cmd)
		// change the container name to IPaddress related
		cmd = "docker"
		args = []string{"rename", cid, strconv.Itoa(count)}
		exe_cmd(cmd, args)

		// IP address +2
		count = count + 2
	}
}

func exe_cmd_full(cmd string) {
	fmt.Println("CMD is : ", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Println("Error to exec CMD", cmd)
	}
	fmt.Println("Output of CMD :", string(out))
}

func exe_cmd(cmd string, args []string) string {
	fmt.Println("command is ", cmd, " ", args)
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Printf("exec cmd error: %s\n", err)
	}
	//fmt.Printf("exec cmd out: %s\n", out)
	s := string(out)

	return s
}
