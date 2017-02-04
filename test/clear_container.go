package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	// first clean the existing containers
	clear_containers()
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
		fmt.Println("cid is: ", cid)
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
	fmt.Println("Output of CMD :", string(out))
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
