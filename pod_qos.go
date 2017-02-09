package main

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"github.com/coreos/etcd/client"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"golang.org/x/net/context"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	pod_dev                        = "eth0"
	node_dev                       = "ens3"    //"eth0"
	br_int                         = "docker0" //"br-int"
	htb_default_classid            = "8001"
	htb_root_handle                = "1:"
	htb_root_classid               = "1:1"
	htb_high_prio                  = "0"
	htb_mid_prio                   = "5"
	htb_low_prio                   = "7"
	all_ip_traffic                 = "0.0.0.0/0"
	classid_max                    = 102
	qos_json_file                  = "qos.json"
	node_default_inbound_bandwidth = "1000" //1000mbps
	pod_default_inbound_min        = "100"  //10mbps
	pod_default_inbound_max        = "500"  //100mbps
)

type qos_para struct {
	NodeIP          string
	PodID           string
	VlanID          string
	VxlanID         string
	PodIP           string
	Action          string
	InBandWidthMin  string
	InBandWidthMax  string
	OutBandWidthMin string
	OutBandWidthMax string
	PodPriority     string
	ClassID         int
}

type container_info struct {
	id        string
	pid       string
	veth_name string
}
type pod_metadata struct {
	cinfo_list []container_info
	classid    int
	pref       string
}

type qosInput []map[string]string

var classid_pool []int

func main() {

	etcd_server := "127.0.0.1:4001"

	load_pod_qos_policy(etcd_server)

}

func load_pod_qos_policy(etcd_server string) map[string]qos_para {

	pod_info_map := map[string]pod_metadata{}
	cid_pid_map := map[string]string{}

	classid_pool = init_classid_pool(classid_max)

	// We can register as many backends that are supported by libkv
	etcd.Register()

	// Initialize a new store with consul
	kv, err := libkv.NewStore(
		store.ETCD,
		[]string{etcd_server},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)

	if err != nil {
		log.Fatal("Cannot create store")
	}

	intf, err := net.InterfaceByName(node_dev)

	if err != nil {
		log.Fatal("Cannot find interface by name " + node_dev)
	}

	mac := intf.HardwareAddr

	key := "/" + string(mac)
	//println("key: ", key)

	stopCh := make(<-chan struct{})
	events, err := kv.Watch(key, stopCh)
	count := 1
	for {
		select {
		case pair := <-events:

			start := time.Now().UnixNano() / 1000000
			pod_qos := parse_qos_info(etcd_server, key)

			t1 := time.Now().UnixNano() / 1000000

			pod_info_map, cid_pid_map := get_pod_info_map(pod_qos, pod_info_map, cid_pid_map)

			t2 := time.Now().UnixNano() / 1000000
			//config pod outbound bandwidth tc qdisc on eth0 in pod
			//set_pod_eth_outbound_bandwidth(pod_qos, pod_info_map)

			t3 := time.Now().UnixNano() / 1000000
			//config pod inbound bandwidth tc qdisc on veth outside
			//set_pod_veth_inbound_bandwidth(pod_qos, pod_info_map)

			t4 := time.Now().UnixNano() / 1000000
			set_br_inbound_bandwidth(br_int, pod_qos, pod_info_map)

			// start to config the Host
			t5 := time.Now().UnixNano() / 1000000
			//Set_vm_outbound_bandwidth(node_dev, pod_qos, pod_info_map)

			t6 := time.Now().UnixNano() / 1000000

			pod_info_map, cid_pid_map = delete_pod_info_map(pod_qos, pod_info_map, cid_pid_map)

			end := time.Now().UnixNano() / 1000000
			fmt.Printf("update pod qos %d: update time %d|%d|%d|%d|%d|%d|%d|%d, value changed on key %s: new value len=%d ... \n", count,
				t1-start, t2-t1, t3-t2, t4-t3, t5-t4, t6-t5, end-t6, end-start, key, len(pair.Value))
			fmt.Printf("time to config pod eth %d \n", t3-t2)
			fmt.Printf("time to config pod veth %d \n", t4-t3)
			fmt.Printf("time to config br %d \n", t5-t4)
			fmt.Printf("time to config VM %d \n", t6-t5)
			count++
		}
	}
}

func init_classid_pool(classid_max int) []int {

	for i := 2; i <= classid_max; i++ {
		classid_pool = append(classid_pool, i)
	}

	return classid_pool
}

func get_classid(classid_pool []int) int {

	len := len(classid_pool)

	if len >= 3 {
		classid := classid_pool[len-1]
		return classid
	} else {
		return 0
	}
}

func dec_classid_pool(classid_pool []int) []int {

	len := len(classid_pool)

	if len >= 3 {
		return classid_pool[:len-1]
	} else {
		return classid_pool
	}
}

func free_classid(classid int) []int {

	classid_pool = append(classid_pool, classid)
	return classid_pool
}

func parse_qos_info(etcd_server string, key string) map[string]qos_para {

	var data qosInput

	proto := "http://" + etcd_server
	cfg := client.Config{
		Endpoints: []string{proto},
		Transport: client.DefaultTransport,
		//set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	kapi := client.NewKeysAPI(c)

	resp, err := kapi.Get(context.Background(), key, nil)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("Current Node Qos: %s: %s\n", resp.Node.Key, resp.Node.Value)
	err = json.Unmarshal([]byte(resp.Node.Value), &data)

	if err != nil {
		log.Fatal(err)
	}

	return load_pod_qos_local(data)
}

func load_pod_qos_local(data qosInput) map[string]qos_para {

	log.Println("Load pod qos local")
	pod_qos := map[string]qos_para{}

	println("Start to retrive qos file...")
	for i := 0; i < len(data); i++ {
		node_id := data[i]["NodeIP"]
		pod_id := data[i]["PodID"]
		vlan_id := data[i]["VlanID"]
		vxlan_id := data[i]["VxlanID"]
		pod_ip := data[i]["PodIP"]
		action := data[i]["Action"]
		inbandwidth_min := data[i]["InBandWidthMin"]
		inbandwidth_max := data[i]["InBandWidthMax"]
		outbandwidth_min := data[i]["OutBandWidthMin"]
		outbandwidth_max := data[i]["OutBandWidthMax"]
		pod_prio := data[i]["PodPriority"]
		classid := 0

		//fmt.Println("pod ip: "+pod_ip)
		pod_qos[pod_ip] = qos_para{node_id, pod_id, vlan_id, vxlan_id,
			pod_ip, action, inbandwidth_min, inbandwidth_max,
			outbandwidth_min, outbandwidth_max, pod_prio, classid}
	}

	//fmt.Print("pod_qos", pod_qos,"\n")
	return pod_qos
}

func get_pod_info_map(pod_qos map[string]qos_para,
	pod_info_map map[string]pod_metadata,
	cid_pid_map map[string]string) (map[string]pod_metadata, map[string]string) {

	println("Start to get pod ip and pid...")

	//get containter id list
	cmd := "docker"
	args := []string{"ps", "-q"}
	ids := exe_cmd(cmd, args)
	//fmt.Print(ids)

	if len(pod_info_map) == 0 {
		println("pod info map is empty, loading pod qos.")
		for _, container_id := range strings.Split(ids, "\n") {

			//println("container_id:"+container_id+".")
			if container_id == "" {
				continue
			}

			//get container pid
			cmd := "docker"
			args := []string{"inspect", "-f", "{{.State.Pid}}", container_id}
			container_pid := strings.Trim(exe_cmd(cmd, args), "\n")
			//fmt.Printf("container id: %s; pid: %s\n", container_id, container_pid)

			//get container ip address,
			cmd = "nsenter"
			args = []string{"-t", container_pid, "-n", "ifconfig", "eth0"}
			ifconfig_output := exe_cmd(cmd, args)

			for _, line := range strings.Split(ifconfig_output, "\n") {

				if strings.Contains(line, "inet addr") {

					container_ip := strings.Split(strings.Split(line, ":")[1], " ")[0]
					//fmt.Printf("container id: %s; pid: %s; ip addr: %s\n", container_id, container_pid, container_ip)

					new_cinfo := container_info{container_id, container_pid, ""}
					//fmt.Print("new_cinfo: ", new_cinfo)

					if pod_meta, ok := pod_info_map[container_ip]; ok {
						list := pod_meta.cinfo_list
						cinfo_list := append(list, new_cinfo)

						new_pod_meta := pod_metadata{cinfo_list, 0, "0"}

						//fmt.Print("new_pod_meta: ",new_pod_meta)
						pod_info_map[container_ip] = new_pod_meta
						cid_pid_map[container_id] = container_pid
					} else {
						cinfo_list := []container_info{new_cinfo}
						pod_meta := pod_metadata{cinfo_list, 0, "0"}
						pod_info_map[container_ip] = pod_meta
						cid_pid_map[container_id] = container_pid
					}
				}
			}
		}
	} else {

		for ip, val := range pod_qos {

			//skip all and default class
			if ip == "all" || ip == "default" {
				continue
			}

			action := val.Action

			switch action {

			case "add", "change", "":

				if _, ok := pod_qos[ip]; ok {
					log.Println("Warning, already have ", ip, " ", pod_qos[ip], " in pod ip pid map.")
					log.Println("pod qos info of ", ip, " is ", pod_qos[ip])

				} else {
					for _, container_id := range strings.Split(ids, "\n") {
						//println("container_id:"+container_id+".")
						if container_id == "" {
							continue
						}

						if _, ok := cid_pid_map[container_id]; ok {

							continue

						} else {

							//get pid and save into pod_info_map
							//get container pid
							cmd := "docker"
							args := []string{"inspect", "-f", "{{.State.Pid}}", container_id}
							container_pid := strings.Trim(exe_cmd(cmd, args), "\n")
							//fmt.Printf("container id: %s; pid: %s\n", container_id, container_pid)

							//get container ip address,
							cmd = "nsenter"
							args = []string{"-t", container_pid, "-n", "ifconfig", "eth0"}
							ifconfig_output := exe_cmd(cmd, args)

							for _, line := range strings.Split(ifconfig_output, "\n") {

								if strings.Contains(line, "inet addr") {
									container_ip := strings.Split(strings.Split(line, ":")[1], " ")[0]
									//fmt.Printf("container id: %s; pid: %s; ip addr: %s\n", container_id, container_pid, container_ip)
									new_cinfo := container_info{container_id, container_pid, ""}

									if metadata, ok := pod_info_map[container_ip]; ok {
										list := metadata.cinfo_list
										cinfo_list := append(list, new_cinfo)

										pod_meta := pod_metadata{cinfo_list, 0, "0"}
										pod_info_map[container_ip] = pod_meta
										cid_pid_map[container_id] = container_pid
									} else {
										cinfo_list := []container_info{new_cinfo}
										pod_meta := pod_metadata{cinfo_list, 0, "0"}
										pod_info_map[container_ip] = pod_meta
										cid_pid_map[container_id] = container_pid
									}
								}
							}
						}
					}
				}

			case "delete":

			default:
			}
		}
	}

	log.Println("pod_info_map: ", pod_info_map)
	log.Println("cid_pid_map: ", cid_pid_map)
	return pod_info_map, cid_pid_map
}

func delete_pod_info_map(pod_qos map[string]qos_para,
	pod_info_map map[string]pod_metadata,
	cid_pid_map map[string]string) (map[string]pod_metadata, map[string]string) {

	//println("Start to delete pod ip and pid...")

	for ip, val := range pod_qos {

		//skip all and default class
		if ip == "all" || ip == "default" {
			continue
		}

		if val.Action == "delete" {
			println("delete", ip)
			//delete container id from cid_pid_map
			cinfo_list := pod_info_map[ip].cinfo_list

			for _, cinfo := range cinfo_list {

				delete(cid_pid_map, cinfo.id)
			}

			//delete container ip from pod_info_map
			delete(pod_info_map, ip)
		}

	}

	//println("Delete pod in pod_info_map")
	//fmt.Print("pod_info_map", pod_info_map,"\n")
	//fmt.Print("cid_pid_map",cid_pid_map, "\n")
	return pod_info_map, cid_pid_map
}

func set_br_inbound_bandwidth(br_name string, pod_qos map[string]qos_para, pod_info_map map[string]pod_metadata) {

	/*tc qdisc add dev $nic root handle 1: htb default 1001
	 *tc class add dev $nic parent 1: classid 1:1 htb rate 10mbit ceil 10mbit
	 *tc class add dev $nic parent 1:1 classid 1:10 htb rate 1mbit ceil 1mbit prio 0
	 *tc class add dev $nic parent 1:1 classid 1:1001 htb rate 8mbit ceil 8mbit prio 3
	 *tc filter add dev $nic parent 1: protocol ip prio 0 u32 match ip dst 0.0.0.0/0 flowid 1:1
	 *tc filter add dev $nic parent 1:1 protocol ip prio 0 u32 match ip dst 10.0.3.153/32 flowid 1:10
	 */

	//println("Start to set bridge inbound bandwidth...")
	//get the sum of pod bandwidth
	node_inbound_bandwidth := node_default_inbound_bandwidth + "mbit"
	node_outbound_bandwidth := node_default_inbound_bandwidth + "mbit"
	action := ""
	ip := "all"
	if _, ok := pod_qos[ip]; ok {
		node_inbound_bandwidth = pod_qos[ip].InBandWidthMax
		node_outbound_bandwidth = pod_qos[ip].OutBandWidthMax
		action = pod_qos[ip].Action
	}

	intf_name := node_dev

	switch action {

	case "add":

		//configure tc qdisc htb on br-int
		//Firstly, delete the tc qdisc on the br-int, tc qdisc del dev br root
		cmd := "tc"
		args := []string{"qdisc", "del", "dev", br_name, "root"}
		exe_cmd(cmd, args)

		//set tc qdisc htb root
		cmd = "tc"
		args = []string{"qdisc", "add", "dev", br_name, "root", "handle", htb_root_handle, "htb", "default", htb_default_classid}
		exe_cmd(cmd, args)

		//set tc class htb 1:1
		//tc class add dev $nic parent 1: classid 1:1 htb rate 10mbit ceil 10mbit
		rate := node_inbound_bandwidth + "mbit"
		cmd = "tc"
		args = []string{"class", "add", "dev", br_name, "parent", htb_root_handle, "classid", htb_root_classid, "htb", "rate", rate, "ceil", rate}
		exe_cmd(cmd, args)

		// Config the VM
		//Firstly, delete the tc qdisc on vm interface, tc qdisc del dev br root
		log.Println(" Set VM interface all")
		cmd = "tc"
		args = []string{"qdisc", "del", "dev", intf_name, "root"}
		exe_cmd(cmd, args)
		//set tc qdisc htb root
		cmd = "tc"
		args = []string{"qdisc", "add", "dev", intf_name, "root", "handle", htb_root_handle, "htb", "default", htb_default_classid}
		exe_cmd(cmd, args)
		//set tc class htb 1:1
		//tc class add dev $nic parent 1: classid 1:1 htb rate 10mbit ceil 10mbit
		rate = node_outbound_bandwidth + "mbit"
		cmd = "tc"
		args = []string{"class", "add", "dev", intf_name, "parent", htb_root_handle, "classid", htb_root_classid, "htb", "rate", rate, "ceil", rate}
		exe_cmd(cmd, args)

	case "delete":

		//  delete the br-int
		println("Delete pod", ip)
		rate := node_inbound_bandwidth + "mbit"
		cmd := "tc"
		args := []string{"class", "del", "dev", br_name, "parent", htb_root_handle, "classid", htb_root_classid, "htb", "rate", rate, "ceil", rate}
		exe_cmd(cmd, args)

		// delete the vm node_dev
		rate = node_outbound_bandwidth + "mbit"
		cmd = "tc"
		args = []string{"class", "del", "dev", intf_name, "parent", htb_root_handle, "classid", htb_root_classid, "htb", "rate", rate, "ceil", rate}
		exe_cmd(cmd, args)

	case "change":
		// change the br-int
		rate := node_inbound_bandwidth + "mbit"
		cmd := "tc"
		args := []string{"class", "change", "dev", br_name, "parent", htb_root_handle, "classid", htb_root_classid, "htb", "rate", rate, "ceil", rate}
		exe_cmd(cmd, args)

		// change the VM node_dev
		rate = node_outbound_bandwidth + "mbit"
		cmd = "tc"
		args = []string{"class", "change", "dev", intf_name, "parent", htb_root_handle, "classid", htb_root_classid, "htb", "rate", rate, "ceil", rate}
		exe_cmd(cmd, args)

	case "":
		println("Not change Qos on pod", ip)

	default:

	}

	//configure default class
	ip = "default"
	action = ""
	rate := pod_default_inbound_min + "mbit"
	ceil := pod_default_inbound_max + "mbit"

	if _, ok := pod_qos[ip]; ok {
		rate = pod_qos[ip].InBandWidthMin + "mbit"
		ceil = pod_qos[ip].InBandWidthMax + "mbit"
		action = pod_qos[ip].Action
	}

	/*
		htb_default_classid = "8001"
		htb_root_handle = "1:"
		htb_root_classid = "1:1"
	*/
	htb_default_classid_full := htb_root_handle + htb_default_classid

	switch action {

	case "add":
		// br-int
		cmd := "tc"
		args := []string{"class", "add", "dev", br_name, "parent", htb_root_classid, "classid", htb_default_classid, "htb", "rate", rate, "ceil", ceil}
		exe_cmd(cmd, args)
		// VM node_dev
		log.Println(" Set VM interface Default")
		cmd = "tc"
		args = []string{"class", "add", "dev", intf_name, "parent", htb_root_classid, "classid", htb_default_classid_full, "htb", "rate", rate, "ceil", ceil}
		exe_cmd(cmd, args)

	case "delete":
		// br-int
		println("Delete pod", ip)
		cmd := "tc"
		args := []string{"class", "del", "dev", br_name, "parent", htb_root_classid, "classid", htb_default_classid, "htb", "rate", rate, "ceil", ceil}
		exe_cmd(cmd, args)
		// vm intface
		cmd = "tc"
		args = []string{"class", "del", "dev", intf_name, "parent", htb_root_classid, "classid", htb_default_classid_full, "htb", "rate", rate, "ceil", ceil}
		exe_cmd(cmd, args)

	case "change":

		cmd := "tc"
		args := []string{"class", "change", "dev", br_name, "parent", htb_root_classid, "classid", htb_default_classid, "htb", "rate", rate, "ceil", ceil}
		exe_cmd(cmd, args)

	case "":
		log.Println("Not change Qos on pod", ip)

	default:

	}

	log.Println("show bridge root and default qdisc and class")
	show_tc_qdisc(br_name)
	show_tc_class(br_name)

	log.Println("show VM root and default qdisc and class")
	show_tc_qdisc(intf_name)
	show_tc_class(intf_name)

	//set tc class and filter for each pod
	set_pod_br_inbound_bandwidth_class_and_filter(br_name, pod_qos, pod_info_map)

	//show tc configuration
	log.Println("show br_int qdisc and class")
	show_tc_qdisc(br_name)
	show_tc_class(br_name)
	log.Println("show br_int filter")
	show_tc_filter(br_name, htb_root_handle)
	show_tc_filter(br_name, htb_root_classid)

	log.Println("show VM qdisc and class")
	show_tc_qdisc(intf_name)
	show_tc_class(intf_name)
	log.Println("show VM filter")
	show_tc_filter(intf_name, htb_root_handle)
	show_tc_filter(intf_name, htb_root_classid)

}

func set_pod_br_inbound_bandwidth_class_and_filter(br_name string, pod_qos map[string]qos_para,
	pod_info_map map[string]pod_metadata) {

	//println("\nStart to set pod inbound bandwidth class and filter on bridge")

	intf_name := node_dev

	for ip, val := range pod_qos {

		//skip all and default class
		if ip == "all" || ip == "default" {
			continue
		}

		rate := val.InBandWidthMin + "mbit"
		ceil := val.InBandWidthMax + "mbit"
		prio := val.PodPriority

		action := val.Action

		switch action {

		case "add":
			/*
				config br-int
			*/
			//configure tc qdisc htb class for each pod on br_int
			classid := get_classid(classid_pool)
			classid_pool = dec_classid_pool(classid_pool)
			cur_classid := htb_root_handle + strconv.Itoa(classid)

			log.Println("Add class and filters with current classID: " + cur_classid)

			//println(ip,action," inbound: "+val.InBandWidthMin+", "+val.InBandWidthMax+", "+val.PodPriority, cur_classid)

			cmd := "tc"
			args := []string{"class", "add", "dev", br_name, "parent", htb_root_classid, "classid",
				cur_classid, "htb", "rate", rate, "ceil", ceil, "prio", prio}
			exe_cmd(cmd, args)

			//set tc filter for each pod on br_int
			//tc filter add dev $nic parent 1:1 protocol ip prio 0 u32 match ip dst 10.0.3.153/32 flowid 1:2
			//println(classid, cur_classid)
			cmd = "tc"
			args = []string{"filter", "add", "dev", br_name, "parent", htb_root_classid, "protocol", "ip",
				"prio", "0", "u32", "match", "ip", "dst", ip + "/32", "flowid", cur_classid}
			exe_cmd(cmd, args)

			/*
				config VM interface
			*/
			cmd = "tc"
			args = []string{"class", "add", "dev", intf_name, "parent", htb_root_classid, "classid",
				cur_classid, "htb", "rate", rate, "ceil", ceil, "prio", prio}
			exe_cmd(cmd, args)

			//filter cmd is "sudo tc filter add dev ens3 parent 1:0 bpf bytecode \"11,40 0 0 12,21 0 8 2048,48 0 0 23,21 0 6 17,40 0 0 42,69 1 0 2048,6 0 0 0,32 0 0 76,21 0 1 167838213,6 0 0 262144,6 0 0 0,\" flowid 1:100"
			// get the byte code
			bytecode := generate_bytecode(ip)
			filterCmd := "tc filter add dev " + string(intf_name) + " parent " + htb_root_classid + " prio " + string(prio) + " bpf bytecode " + bytecode + " flowid " + cur_classid
			exe_cmd_full(filterCmd)

			//get filter pref,
			cmd = "tc"
			args = []string{"filter", "show", "dev", br_name, "parent", htb_root_classid}
			output := exe_cmd(cmd, args)

			var pref string
			for _, line := range strings.Split(output, "\n") {

				if strings.Contains(line, cur_classid) {
					pref = strings.Split(line, " ")[4]
				}
			}

			//update classid in pod_info_map
			if _, ok := pod_info_map[ip]; ok {
				pod_meta := pod_info_map[ip]
				pod_meta.classid = classid
				pod_meta.pref = pref
				//fmt.Print(pod_meta)
				delete(pod_info_map, ip)
				pod_info_map[ip] = pod_meta

				//fmt.Print("pod info:", pod_info_map[ip])
			} else {
				println("can not find in pod info map.", ip)
			}

		case "delete":

			if _, ok := pod_info_map[ip]; ok {

				classid := pod_info_map[ip].classid
				cur_classid := htb_root_handle + strconv.Itoa(classid)
				pref := pod_info_map[ip].pref

				//println("Delete pod filter on",cur_classid, br_name, ip, pref)

				cmd := "tc"
				args := []string{"filter", "del", "dev", br_name, "parent", htb_root_classid, "prio", pref, "u32"}

				exe_cmd(cmd, args)

				//println("Delete pod class on",cur_classid, br_name, ip, pref)
				cmd = "tc"
				args = []string{"class", "del", "dev", br_name, "parent", htb_root_classid, "classid",
					cur_classid, "htb", "rate", rate, "ceil", ceil, "prio", prio}
				exe_cmd(cmd, args)

				/*
					delete VM interface
					sudo tc filter del dev ens3 parent 1: prio 1
					sudo tc class del dev ens3 parent 1:1 classid 1:95
					note: there is a different with br-int that the prio and pref
				*/
				log.Println("delete filter on ", intf_name)
				/*
					when prio = 0, the pref in show filter is 49152
				*/
				if prio == "0" {
					log.Println("change prio 0 to pref 49152")
					prio = "49152"
				}
				cmd = "tc"
				args = []string{"filter", "del", "dev", intf_name, "parent", htb_root_classid, "prio", prio}
				exe_cmd(cmd, args)

				log.Println("delete class on ", intf_name)
				cmd = "tc"
				args = []string{"class", "del", "dev", intf_name, "parent", htb_root_classid, "classid",
					cur_classid}
				exe_cmd(cmd, args)

				// update classid_pool
				classid_pool = free_classid(classid)

				pod_meta := pod_info_map[ip]
				pod_meta.classid = 0
				pod_meta.pref = "0"
				pod_info_map[ip] = pod_meta
				//fmt.Print(pod_info_map[ip])
			} else {
				log.Println("Can NOT find the config for " + ip + " in pod info map.")
			}

		case "change":

			classid := pod_info_map[ip].classid
			cur_classid := htb_root_handle + strconv.Itoa(classid)
			log.Println("change class on" + br_name + " and current classID is: " + cur_classid)
			/*
				config br-int
			*/
			cmd := "tc"
			args := []string{"class", "change", "dev", br_name, "parent", htb_root_classid, "classid",
				cur_classid, "htb", "rate", rate, "ceil", ceil, "prio", prio}
			exe_cmd(cmd, args)

			/*
				config br-int
			*/
			cmd = "tc"
			args = []string{"class", "change", "dev", intf_name, "parent", htb_root_classid, "classid",
				cur_classid, "htb", "rate", rate, "ceil", ceil, "prio", prio}
			exe_cmd(cmd, args)

		case "":
			println("Not change Qos on pod", ip)

		default:

		}

		if len(classid_pool) == 0 {
			println("Error classid pool is empty. Cannot set pod bridge inbound bandwidth class and filter.")
			break
		}
	}

	//return classid_pool
}

func generate_bytecode(ip string) string {
	//log.Println("current IP is : ", ip)
	// for Vxlan src 10.1.2.5
	//sudo tc filter add dev ens3 parent 1:0 bpf bytecode \
	//"11,40 0 0 12,21 0 8 2048,48 0 0 23,21 0 6 17,40 0 0 42,69 1 0 2048,6 0 0 0,32 0 0 76,21 0 1 167838213,6 0 0 262144,6 0 0 0," flowi d 1:20
	part1 := "\"11,40 0 0 12,21 0 8 2048,48 0 0 23,21 0 6 17,40 0 0 42,69 1 0 2048,6 0 0 0,32 0 0 76,21 0 1 "
	var temp int64
	for _, value := range strings.Split(ip, ".") {
		temp = temp << 8
		i, err := strconv.Atoi(value)
		//fmt.Println("value: ", i)
		if err != nil {
			log.Fatal("Error in Atoi ")
		}
		temp = temp + int64(i)
	}
	part2 := strconv.FormatInt(temp, 10)
	part3 := ",6 0 0 262144,6 0 0 0,\""
	code := part1 + part2 + part3
	//log.Println("bytecode is: ", code)
	return code
}

func set_pod_veth_inbound_bandwidth(pod_qos map[string]qos_para, pod_info_map map[string]pod_metadata) {

	//config pod inbound bandwidth tc qdisc on veth outside
	veth_list := get_veth_list()

	for ip, val := range pod_qos {
		//println(ip+" inbound: "+val.InBandWidthMin+", "+val.InBandWidthMax)

		//if ip = "all", skip it
		if ip == "all" || ip == "default" {
			//println("all ip traffic bandwidth.")
			continue
		}

		//get container pid
		if pod_meta, ok := pod_info_map[ip]; ok {
			container_pid := pod_meta.cinfo_list[0].pid
			action := val.Action

			switch action {

			case "add":

				//get veth id
				cmd := "nsenter"
				args := []string{"-t", container_pid, "-n", "ip", "link", "show", "eth0"}
				output := exe_cmd(cmd, args)

				dev_id := strings.Split(output, ":")[0]
				tmp, _ := strconv.Atoi(dev_id)
				veth_id := tmp + 1
				//println("dev_id:",tmp, ", veth id:",strconv.Itoa(veth_id))

				veth_name := veth_list[veth_id]
				//println(strconv.Itoa(veth_id), veth_name)

				//configure tc qdisc tbf on eth0 in the container(pod)
				//Firstly, delete the tc qdisc on the veth, tc qdisc del dev veth_id root
				cmd = "tc"
				args = []string{"qdisc", "del", "dev", veth_name, "root"}
				exe_cmd(cmd, args)

				//secondly, configure tc qdisc tbf
				//tc qdisc add dev eth0 root tbf rate mbit latency 50ms burst 100k
				cmd = "tc"
				args = []string{"qdisc", "add", "dev", veth_name, "root", "tbf", "rate",
					val.InBandWidthMax + "mbit", "latency", "50ms", "burst", "100k"}
				exe_cmd(cmd, args)

				//show tc configuration
				println("add qos on ", ip, veth_name)
				//show_tc_qdisc(veth_name)

			case "delete":

				//get veth id
				cmd := "nsenter"
				args := []string{"-t", container_pid, "-n", "ip", "link", "show", "eth0"}
				output := exe_cmd(cmd, args)

				dev_id := strings.Split(output, ":")[0]
				tmp, _ := strconv.Atoi(dev_id)
				veth_id := tmp + 1
				//println("dev_id:",tmp, ", veth id:",strconv.Itoa(veth_id))

				veth_name := veth_list[veth_id]
				//println(strconv.Itoa(veth_id), veth_name)

				//configure tc qdisc tbf on eth0 in the container(pod)
				//Firstly, delete the tc qdisc on the veth, tc qdisc del dev veth_id root
				cmd = "tc"
				args = []string{"qdisc", "del", "dev", veth_name, "root"}
				exe_cmd(cmd, args)

				//show tc configuration
				println("delete qos on ", ip, veth_name)
				//show_tc_qdisc(veth_name)

			case "change":

				//get veth id
				cmd := "nsenter"
				args := []string{"-t", container_pid, "-n", "ip", "link", "show", "eth0"}
				output := exe_cmd(cmd, args)

				dev_id := strings.Split(output, ":")[0]
				tmp, _ := strconv.Atoi(dev_id)
				veth_id := tmp + 1
				//println("dev_id:",tmp, ", veth id:",strconv.Itoa(veth_id))

				veth_name := veth_list[veth_id]
				println("Change pod Qos on", ip, veth_name)
				cmd = "tc"
				args = []string{"qdisc", "change", "dev", veth_name, "root", "tbf", "rate",
					val.InBandWidthMax + "mbit", "latency", "50ms", "burst", "100k"}
				exe_cmd(cmd, args)

				//show tc configuration
				println("change qos on ", ip, veth_name)
				//show_tc_qdisc(veth_name)

			case "":
				println("Not change Qos on pod", ip)

			default:

			}

		} else {
			println("Can not find ip: " + ip + " in pod_info_map.")
			continue
		}
	}
}

func set_pod_eth_outbound_bandwidth(pod_qos map[string]qos_para, pod_info_map map[string]pod_metadata) {

	//config pod outbound bandwidth tc qdisc on eth0 in pod
	//println("Start to set pod outbound bandwidth in pod...")

	for ip, val := range pod_qos {
		println(ip + " outbound: " + val.OutBandWidthMin + ", " + val.OutBandWidthMax)

		//if ip = "all", skip it
		if ip == "all" || ip == "default" {
			//println("all ip traffic bandwidth.")
			continue
		}
		//get container pid
		if pod_meta, ok := pod_info_map[ip]; ok {

			container_pid := pod_meta.cinfo_list[0].pid
			action := val.Action

			switch action {

			case "add":

				//configure tc qdisc tbf on eth0 in the container(pod)
				//Firstly, delete the tc qdisc on the eth0, tc qdisc del dev $1 root
				cmd := "nsenter"
				args := []string{"-t", container_pid, "-n", "tc", "qdisc", "del", "dev", pod_dev, "root"}
				exe_cmd(cmd, args)

				//secondly, configure tc qdisc tbf
				//tc qdisc add dev eth0 root tbf rate mbit latency 50ms burst 100k
				cmd = "nsenter"
				args = []string{"-t", container_pid, "-n", "tc", "qdisc", "add", "dev", "eth0", "root", "tbf", "rate",
					val.OutBandWidthMax + "mbit", "latency", "50ms", "burst", "100k"}
				exe_cmd(cmd, args)
				println("Add qos on pod", ip, pod_dev)
				show_tc_qdisc_in_pod(container_pid, pod_dev)

			case "delete":

				println("Delete pod Qos on", ip, pod_dev)
				cmd := "nsenter"
				args := []string{"-t", container_pid, "-n", "tc", "qdisc", "del", "dev", pod_dev, "root"}
				exe_cmd(cmd, args)
				show_tc_qdisc_in_pod(container_pid, pod_dev)

			case "change":

				println("Change pod Qos on", ip, pod_dev)

				cmd := "nsenter"
				args := []string{"-t", container_pid, "-n", "tc", "qdisc", "change", "dev", "eth0", "root", "tbf", "rate",
					val.OutBandWidthMax + "mbit", "latency", "50ms", "burst", "100k"}
				exe_cmd(cmd, args)
				show_tc_qdisc_in_pod(container_pid, pod_dev)

			case "":
				log.Println("Not change Qos on pod", ip, pod_dev)

			default:

			}
			//show tc configuration
			//println("pod "+ip+" tc qdisc show: ")
			//show_tc_qdisc_in_pod(container_pid, pod_dev)

		} else {
			println("Can not find ip: " + ip + " in pod_info_map.")
			continue
		}
	}

}

//get veth list on the host
func get_veth_list() map[int]string {

	//println("Start to get veth list...")
	result := map[int]string{}
	intf_list, err := net.Interfaces()
	if err != nil {
		println(err)

	}
	for _, f := range intf_list {

		if strings.Contains(f.Name, "veth") {
			result[f.Index] = f.Name
			//fmt.Println(f.Index, f.Name)
		}
	}
	return result
}

//get interface and it's IP address on VM
func get_intf_ipaddress(intf_name string) net.IP {
	var result net.IP
	result = nil
	ifaces, err := net.Interfaces()

	if err != nil {
		fmt.Printf("Error when decode interface %s\n", err)
	}
	// handle err
	for _, i := range ifaces {
		if i.Name == intf_name {
			addrs, err := i.Addrs()
			// handle err
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
					fmt.Println("IP net is: ", ip)
					return ip
				case *net.IPAddr:
					//result = v.IP
					ip = v.IP
					fmt.Println("IP address is: ", ip)
				}
				// process IP address
			}

			if err != nil {
				fmt.Printf("Error when decode interface %s\n", err)
			}
		}
	}
	return result
}

func exe_cmd_full(cmd string) {
	log.Println("command is : ", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	//_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("Error to exec CMD", cmd)
	}
	log.Println("Output of command:", string(out))
}

func exe_cmd(cmd string, args []string) string {

	log.Println("command is ", cmd, " ", args)

	out, err := exec.Command(cmd, args...).Output()

	if err != nil {
		fmt.Printf("exec cmd error: %s\n", err)
	}

	//fmt.Printf("exec cmd out: %s\n", out)

	s := string(out)

	return s
}

func show_tc_qdisc(dev_name string) {

	//show tc configuration
	cmd := "tc"
	args := []string{"qdisc", "show", "dev", dev_name}
	println(exe_cmd(cmd, args))

}

func show_tc_class(dev_name string) {

	//show tc configuration
	cmd := "tc"
	args := []string{"class", "show", "dev", dev_name}
	println(exe_cmd(cmd, args))

}

func show_tc_qdisc_statistics(dev_name string) {

	cmd := "tc"
	args := []string{"-s", "qdisc", "show", "dev", dev_name}
	println(exe_cmd(cmd, args))

}

func show_tc_class_statistics(dev_name string) {

	cmd := "tc"
	args := []string{"-s", "class", "show", "dev", dev_name}
	println(exe_cmd(cmd, args))

}

func show_tc_qdisc_in_pod(container_pid string, dev_name string) {

	//show tc configuration
	cmd := "nsenter"
	args := []string{"-t", container_pid, "-n", "tc", "qdisc", "show", "dev", dev_name}
	println(exe_cmd(cmd, args))

}

func show_tc_class_in_pod(container_pid string, dev_name string) {

	//show tc configuration
	cmd := "nsenter"
	args := []string{"-t", container_pid, "-n", "tc", "class", "show", "dev", dev_name}
	println(exe_cmd(cmd, args))

}

func show_tc_qdisc_statistics_in_pod(container_pid string, dev_name string) {

	//show tc configuration
	cmd := "nsenter"
	args := []string{"-t", container_pid, "-n", "tc", "-s", "qdisc", "show", "dev", dev_name}
	println(exe_cmd(cmd, args))

}

func show_tc_class_statistics_in_pod(container_pid string, dev_name string) {

	//show tc configuration
	cmd := "nsenter"
	args := []string{"-t", container_pid, "-n", "tc", "-s", "class", "show", "dev", dev_name}
	println(exe_cmd(cmd, args))

}

func show_tc_filter(dev_name string, handle string) {

	//show tc configuration
	cmd := "tc"
	args := []string{"filter", "show", "dev", dev_name, "parent", handle}
	println(exe_cmd(cmd, args))

}
