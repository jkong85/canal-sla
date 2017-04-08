package main

import (
	"common"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

var podDev = "eth0"
var nodeDev = "enp0s9"
var bridgeDev = "br-int"

const (
	htbDefaultClassid           = "8001"
	htbRootHandle               = "1:"
	htbRootClassid              = "1:1"
	htbHighPrio                 = "0"
	htbMidPrio                  = "5"
	htbLowPrio                  = "7"
	allIPTraffic                = "0.0.0.0/0"
	classidMax                  = 102
	qosJSONFile                 = "qos.json"
	nodeDefaultInboundBandwidth = "1000" //1000mbps
	podDefaultInboundMin        = "100"  //10mbps
	podDefaultInboundMax        = "500"  //100mbps
)

type qosPara struct {
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

type containerInfo struct {
	id       string
	pid      string
	vethName string
}
type podMetadata struct {
	cinfoList []containerInfo
	classid   int
	pref      string
}
type qosInput []map[string]string

var classidPool []int

/*
func main() {
	etcd_server := "127.0.0.1:4001"
	load_podQos_policy(etcd_server)
}
*/

var count = 0

// LoadPodQosPolicy : Configure the Qos policy for the interfaces
func LoadPodQosPolicy(val string, podname string, ip string, action string) {
	if strings.Compare(val, "") == 0 {
		return
	}

	podInfoMap := map[string]podMetadata{}
	cidPidMap := map[string]string{}
	//podDev, _ := GetProviderInterface()
	//nodeDev = eth.Name
	nodeDev, _, _ = common.GetInterface("default")

	bridgeDev = getBridgeName()

	if count == 0 { // initialize when the first call
		log.Println("Initialize when the first call!")
		classidPool = initClassidPool(classidMax)
	}

	//get the IP address of this host
	log.Printf("Provider interface %s\n", nodeDev)

	//start := time.Now().UnixNano() / 1000000
	podQos, err := parseQosInfo(val, podname, ip, action)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Start to Load Qos Policy : %d", count)
	log.Println("Before load the configuration, podInfoMap is", podInfoMap)

	// set the All and default ONCE (the first time trigged)
	if count == 0 {
		log.Println("set up the default and all")
		setHTBAllDefault(bridgeDev)
	}

	podInfoMap, cidPidMap = getPodInfoMap(podQos, podInfoMap, cidPidMap)
	log.Println("Before device config, the podInfoMap is", podInfoMap)

	setPodTC(podQos, podInfoMap)

	setVethTC(podQos, podInfoMap)

	//setBrVmTC(bridgeDev, podQos, podInfoMap)
	setBrVmClassFilter(bridgeDev, podQos, podInfoMap)

	// update the podInfoMap ( if the action is delete, delete from the podInfoMap)
	podInfoMap, cidPidMap = updatePodInfoMap(podQos, podInfoMap, cidPidMap)
	log.Println("update the PodInfoMap : ", podInfoMap)
	log.Println("Load QOS is finished !")

	count++
}

func parseQosInfo(val string, podname string, ip string, action string) (map[string]qosPara, error) {
	var data qosInput
	err := json.Unmarshal([]byte(val), &data)
	if err != nil {
		log.Fatal(err)
	}
	return loadPodQosLocal(data, podname, ip, action)
}

func loadPodQosLocal(data qosInput, podname string, ip string, action string) (map[string]qosPara, error) {

	log.Println("Load pod qos local")
	podQos := map[string]qosPara{}

	// check whether "All" action is delete, if yes, it means that we need to delete all rules from the Pod, veth, br-int and VM interface. we set all rules's action to 'delete' to realize it
	flagDeleteAll := false

	for i := 0; i < len(data); i++ {
		if data[i]["PodIP"] == "all" && action == "delete" {
			log.Println("delete All")
			flagDeleteAll = true
			break
		}
	}
	var i int
	for i = 0; i < len(data); i++ {
		nodeID := data[i]["NodeIP"]
		podID := data[i]["PodID"]
		vlanID := data[i]["VlanID"]
		vxlanID := data[i]["VxlanID"]
		podIP := data[i]["PodIP"]

		if strings.Compare(podID, podname) != 0 {
			continue
		}

		//use podIP from topo listener
		podIP = ip

		//action := data[i]["Action"]
		inBandwidthMin := data[i]["InBandWidthMin"]
		inBandwidthMax := data[i]["InBandWidthMax"]
		outBandwidthMin := data[i]["OutBandWidthMin"]
		outBandwidthMax := data[i]["OutBandWidthMax"]
		podPriority := data[i]["PodPriority"]
		classid := 0

		// If the rule is for current host, then add it to the podQos
		if flagDeleteAll {
			// set all rules action to "delete"
			log.Println("set all rules' action to delete")
			action = "delete"
		}
		podQos[podIP] = qosPara{nodeID, podID, vlanID, vxlanID,
			podIP, action, inBandwidthMin, inBandwidthMax,
			outBandwidthMin, outBandwidthMax, podPriority, classid}

		break
	}

	if i == len(data) {
		//log.Fatalf("App not in the nqos map")
		log.Printf("App %s is not in nqos map %+v\n", podname, data)
		return podQos, errors.New("App not in the nqos map")
	}
	log.Println("podQos after load is: ", podQos)
	return podQos, nil
}

// Input: podQos, actionType: 0 -> change all to delete, 1 -> change " ", "change" to "add"
func changeAction(podQosOld map[string]qosPara, actionType int) map[string]qosPara {
	podQosChange := map[string]qosPara{}

	for _, val := range podQosOld {
		nodeID := val.NodeIP
		podID := val.PodID
		vlanID := val.VlanID
		vxlanID := val.VxlanID
		podIP := val.PodIP
		action := val.Action
		inBandwidthMin := val.InBandWidthMin
		inBandwidthMax := val.InBandWidthMax
		outBandwidthMin := val.OutBandWidthMin
		outBandwidthMax := val.OutBandWidthMax
		podPriority := val.PodPriority
		classid := val.ClassID

		if actionType == 0 {
			action = "delete"
		} else if actionType == 1 {
			switch action {
			case "", " ", "change":
				action = "add"
			case "add":
			default:

			}
		} else {
			log.Println("Specify the correct action!")
		}

		podQosChange[podIP] = qosPara{nodeID, podID, vlanID, vxlanID,
			podIP, action, inBandwidthMin, inBandwidthMax,
			outBandwidthMin, outBandwidthMax, podPriority, classid}

	}
	//log.Println("After change, the podQos is: ", podQosChange)
	return podQosChange
}

func getPodInfoMap(podQos map[string]qosPara,
	podInfoMap map[string]podMetadata,
	cidPidMap map[string]string) (map[string]podMetadata, map[string]string) {

	log.Println("Start to get pod info ")

	//get containter id list
	cmd := "docker ps -q"
	ids := exeCMDFull(cmd)

	if len(podInfoMap) == 0 {
		log.Println("pod info map is empty, loading pod qos.")
		for _, containerID := range strings.Split(ids, "\n") {

			//log.Println("containerID:"+containerID+".")
			if containerID == "" {
				continue
			}

			//get container pid
			//cmd := "docker"
			//args := []string{"inspect", "-f", "{{.State.Pid}}", containerID}
			cmd := "docker inspect -f {{.State.Pid}} " + containerID
			containerPID := strings.Trim(exeCMDFull(cmd), "\n")
			//fmt.Printf("container id: %s; pid: %s\n", containerID, containerPID)

			//get container ip address,
			//cmd = "nsenter"
			//args = []string{"-t", containerPID, "-n", "ifconfig", "eth0"}
			cmd = "nsenter -t " + containerPID + " -n ifconfig eth0 "
			ifconfigOutput := exeCMDFull(cmd)

			for _, line := range strings.Split(ifconfigOutput, "\n") {

				if strings.Contains(line, "inet addr") {

					containerIP := strings.Split(strings.Split(line, ":")[1], " ")[0]
					//fmt.Printf("container id: %s; pid: %s; ip addr: %s\n", containerID, containerPID, containerIP)

					newConInfo := containerInfo{containerID, containerPID, ""}
					//fmt.Print("newConInfo: ", newConInfo)

					if podMeta, ok := podInfoMap[containerIP]; ok {
						list := podMeta.cinfoList
						cinfoList := append(list, newConInfo)

						newPodMeta := podMetadata{cinfoList, 0, "0"}

						//fmt.Print("newPodMeta: ",newPodMeta)
						podInfoMap[containerIP] = newPodMeta
						cidPidMap[containerID] = containerPID
					} else {
						cinfoList := []containerInfo{newConInfo}
						podMeta := podMetadata{cinfoList, 0, "0"}
						podInfoMap[containerIP] = podMeta
						cidPidMap[containerID] = containerPID
					}
				}
			}
		}
	} else {
		for ip, val := range podQos {
			//skip all and default class
			if ip == "all" || ip == "default" {
				continue
			}

			action := val.Action

			switch action {
			case "add", "change", "":
				if _, ok := podQos[ip]; ok {
					log.Println("Warning, already have ", ip, " ", podQos[ip], " in pod ip pid map.")
					log.Println("pod qos info of ", ip, " is ", podQos[ip])
					log.Println("pod info map is :", podInfoMap)
				} else {
					log.Println("There is no such Pod IP in podInfoMap")
					for _, containerID := range strings.Split(ids, "\n") {
						//log.Println("containerID:"+containerID+".")
						if containerID == "" {
							// there is no pod existing for the qos policy, ignore it
							continue
						}

						if _, ok := cidPidMap[containerID]; ok {
							continue
						} else {
							//get pid and save into podInfoMap
							//get container pid
							//cmd := "docker"
							//args := []string{"inspect", "-f", "{{.State.Pid}}", containerID}
							cmd := "docker inspect -f {{.State.Pid}} " + containerID
							containerPID := strings.Trim(exeCMDFull(cmd), "\n")
							//fmt.Printf("container id: %s; pid: %s\n", containerID, containerPID)

							//get container ip address,
							//cmd = "nsenter"
							//args = []string{"-t", containerPID, "-n", "ifconfig", "eth0"}
							//ifconfigOutput := exeCMD(cmd, args)
							cmd = " nsenter -t " + containerPID + " -n ifconfig eth0 "
							ifconfigOutput := exeCMDFull(cmd)

							for _, line := range strings.Split(ifconfigOutput, "\n") {
								if strings.Contains(line, "inet addr") {
									containerIP := strings.Split(strings.Split(line, ":")[1], " ")[0]
									//fmt.Printf("container id: %s; pid: %s; ip addr: %s\n", containerID, containerPID, containerIP)
									newConInfo := containerInfo{containerID, containerPID, ""}

									if metadata, ok := podInfoMap[containerIP]; ok {
										list := metadata.cinfoList
										cinfoList := append(list, newConInfo)

										podMeta := podMetadata{cinfoList, 0, "0"}
										podInfoMap[containerIP] = podMeta
										cidPidMap[containerID] = containerPID
									} else {
										cinfoList := []containerInfo{newConInfo}
										podMeta := podMetadata{cinfoList, 0, "0"}
										podInfoMap[containerIP] = podMeta
										cidPidMap[containerID] = containerPID
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

	log.Println("podInfoMap: ", podInfoMap)
	log.Println("cidPidMap: ", cidPidMap)
	return podInfoMap, cidPidMap
}

// if the action is delete, then update the podInfoMap
func updatePodInfoMap(podQos map[string]qosPara,
	podInfoMap map[string]podMetadata,
	cidPidMap map[string]string) (map[string]podMetadata, map[string]string) {

	for ip, val := range podQos {
		//skip all and default class
		if ip == "all" || ip == "default" {
			continue
		}
		if val.Action == "delete" {
			log.Println("update podInfoMap by deleting Pod : ", ip)
			//delete container id from cidPidMap
			cinfoList := podInfoMap[ip].cinfoList

			for _, cinfo := range cinfoList {

				delete(cidPidMap, cinfo.id)
			}
			//delete container ip from podInfoMap
			delete(podInfoMap, ip)
		}

	}

	return podInfoMap, cidPidMap
}

func setHTBAllDefault(brName string) {
	/*
		msg.NetRxBw, _ = strconv.ParseUint(common.InBw, 10, 64)
		msg.NetTxBw, _ = strconv.ParseUint(common.OutBw, 10, 64)
	*/
	intfName := nodeDev
	//RxBw := strconv.ParseUint(common.InBw, 10, 64)
	ceil := common.InBw + "mbit"
	rate := common.InBw + "mbit"
	log.Println("The ceiling rate of the All configuration is: " + ceil)
	//configure tc qdisc htb on br-int
	//Firstly, delete the tc qdisc on the br-int, tc qdisc del dev br root
	//cmd := "tc"
	//args := []string{"qdisc", "del", "dev", brName, "root"}
	//exeCMD(cmd, args)
	cmd := " tc qdisc del dev " + brName + " root "
	exeCMDFull(cmd)

	//set tc qdisc htb root
	//cmd = "tc"
	//args = []string{"qdisc", "add", "dev", brName, "root", "handle", htbRootHandle, "htb", "default", htbDefaultClassid}
	//exeCMD(cmd, args)
	cmd = " tc qdisc add dev " + brName + " root handle " + htbRootHandle + " htb default " + htbDefaultClassid
	exeCMDFull(cmd)

	//set tc class htb 1:1
	//tc class add dev $nic parent 1: classid 1:1 htb rate 10mbit ceil 10mbit
	//cmd = "tc"
	//args = []string{"class", "add", "dev", brName, "parent", htbRootHandle, "classid", htbRootClassid, "htb", "rate", rate, "ceil", rate}
	//exeCMD(cmd, args)
	cmd = " tc class add dev " + brName + " parent " + htbRootHandle + " classid " + htbRootClassid + " htb rate " + rate + " ceil " + rate
	exeCMDFull(cmd)

	// Config the VM
	//Firstly, delete the tc qdisc on vm interface, tc qdisc del dev br root
	log.Println(" set the qdisc (All) for the VM interface ")
	//cmd = "tc"
	//args = []string{"qdisc", "del", "dev", intfName, "root"}
	//exeCMD(cmd, args)
	cmd = " tc qdisc del dev " + intfName + " root "
	exeCMDFull(cmd)
	//set tc qdisc htb root
	//cmd = "tc"
	//args = []string{"qdisc", "add", "dev", intfName, "root", "handle", htbRootHandle, "htb", "default", htbDefaultClassid}
	//exeCMD(cmd, args)
	cmd = " tc qdisc add dev " + intfName + " root handle " + htbRootHandle + " htb default " + htbDefaultClassid
	exeCMDFull(cmd)
	//set tc class htb 1:1
	//tc class add dev $nic parent 1: classid 1:1 htb rate 10mbit ceil 10mbit
	//cmd = "tc"
	//args = []string{"class", "add", "dev", intfName, "parent", htbRootHandle, "classid", htbRootClassid, "htb", "rate", rate, "ceil", rate}
	//exeCMD(cmd, args)
	cmd = " tc class add dev " + intfName + " parent " + htbRootHandle + " classid " + htbRootClassid + " htb rate " + rate + " ceil " + rate
	exeCMDFull(cmd)

	//configure default class
	allRateStr, _ := strconv.ParseUint(common.InBw, 10, 64)
	allCeilStr, _ := strconv.ParseUint(common.OutBw, 10, 64)
	defaultRate := strconv.Itoa(int(allRateStr) / 10)
	defaultCeil := strconv.Itoa(int(allCeilStr) / 10)
	rate = defaultRate + "mbit"
	ceil = defaultCeil + "mbit"
	log.Println("The default rate is : " + rate + " and the ceil rate is: " + ceil)

	htbDefaultClassidFull := htbRootHandle + htbDefaultClassid

	// br-int
	//cmd = "tc"
	//args = []string{"class", "add", "dev", brName, "parent", htbRootClassid, "classid", htbDefaultClassid, "htb", "rate", rate, "ceil", ceil}
	//exeCMD(cmd, args)
	log.Println(" Set the default for : " + brName)
	cmd = " tc class add dev " + brName + " parent " + htbRootClassid + " classid " + htbDefaultClassid + " htb rate " + rate + " ceil " + ceil
	exeCMDFull(cmd)
	// VM nodeDev
	//cmd = "tc"
	//args = []string{"class", "add", "dev", intfName, "parent", htbRootClassid, "classid", htbDefaultClassidFull, "htb", "rate", rate, "ceil", ceil}
	//exeCMD(cmd, args)
	log.Println(" Set the default for : " + intfName)
	cmd = " tc class add dev " + intfName + " parent " + htbRootClassid + " classid " + htbDefaultClassidFull + " htb rate " + rate + " ceil " + ceil
	exeCMDFull(cmd)
}

// setBrVmTC : update the VM's interface and the bridge interface at the same time
func setBrVmTC(brName string, podQos map[string]qosPara, podInfoMap map[string]podMetadata) {

	/*tc qdisc add dev $nic root handle 1: htb default 1001
	 *tc class add dev $nic parent 1: classid 1:1 htb rate 10mbit ceil 10mbit
	 *tc class add dev $nic parent 1:1 classid 1:10 htb rate 1mbit ceil 1mbit prio 0
	 *tc class add dev $nic parent 1:1 classid 1:1001 htb rate 8mbit ceil 8mbit prio 3
	 *tc filter add dev $nic parent 1: protocol ip prio 0 u32 match ip dst 0.0.0.0/0 flowid 1:1
	 *tc filter add dev $nic parent 1:1 protocol ip prio 0 u32 match ip dst 10.0.3.153/32 flowid 1:10
	 */
	//get the sum of pod bandwidth
	nodeInbandwidth := nodeDefaultInboundBandwidth + "mbit"
	nodeOutbandwidth := nodeDefaultInboundBandwidth + "mbit"
	action := ""
	ip := "all"
	if _, ok := podQos[ip]; ok {
		nodeInbandwidth = podQos[ip].InBandWidthMax
		nodeOutbandwidth = podQos[ip].OutBandWidthMax
		action = podQos[ip].Action
	}

	intfName := nodeDev

	switch action {
	case "add":
		//configure tc qdisc htb on br-int
		//Firstly, delete the tc qdisc on the br-int, tc qdisc del dev br root
		//cmd := "tc"
		//args := []string{"qdisc", "del", "dev", brName, "root"}
		//exeCMD(cmd, args)
		//cmd := " tc qdisc del dev " + brName + " root "
		//exeCMDFull(cmd)
		args := []string{"qdisc", "del", "dev", brName, "root"}

		//set tc qdisc htb root
		//cmd = "tc"
		//args = []string{"qdisc", "add", "dev", brName, "root", "handle", htbRootHandle, "htb", "default", htbDefaultClassid}
		//exeCMD(cmd, args)
		cmd := " tc qdisc add dev " + brName + " root handle " + htbRootHandle + " htb default " + htbDefaultClassid
		exeCMDFull(cmd)

		//set tc class htb 1:1
		//tc class add dev $nic parent 1: classid 1:1 htb rate 10mbit ceil 10mbit
		rate := nodeInbandwidth + "mbit"
		cmd = "tc"
		args = []string{"class", "add", "dev", brName, "parent", htbRootHandle, "classid", htbRootClassid, "htb", "rate", rate, "ceil", rate}
		exeCMD(cmd, args)

		// Config the VM
		//Firstly, delete the tc qdisc on vm interface, tc qdisc del dev br root
		log.Println(" Set VM interface all")
		//cmd = "tc"
		//args = []string{"qdisc", "del", "dev", intfName, "root"}
		//exeCMD(cmd, args)

		//set tc qdisc htb root
		cmd = "tc"
		args = []string{"qdisc", "add", "dev", intfName, "root", "handle", htbRootHandle, "htb", "default", htbDefaultClassid}
		exeCMD(cmd, args)
		//set tc class htb 1:1
		//tc class add dev $nic parent 1: classid 1:1 htb rate 10mbit ceil 10mbit
		rate = nodeOutbandwidth + "mbit"
		cmd = "tc"
		args = []string{"class", "add", "dev", intfName, "parent", htbRootHandle, "classid", htbRootClassid, "htb", "rate", rate, "ceil", rate}
		exeCMD(cmd, args)

	case "delete":
		// if we delete "ALL", it means we remove all the qdisc config on the interface
		// therefore, we run 'tc qdisc del dev intf root'

		// delete the br-int
		cmd := "tc"
		args := []string{"qdisc", "del", "dev", brName, "root"}
		exeCMD(cmd, args)

		// delete the vm interface
		cmd = "tc"
		args = []string{"qdisc", "del", "dev", intfName, "root"}
		exeCMD(cmd, args)
		/*
			// Original code to delete:
			// delete the br-int
			rate := nodeInbandwidth + "mbit"
			cmd := "tc"
			args := []string{"class", "del", "dev", brName, "parent", htbRootHandle, "classid", htbRootClassid, "htb", "rate", rate, "ceil", rate}
			exeCMD(cmd, args)

			// delete the vm nodeDev
			rate = nodeOutbandwidth + "mbit"
			cmd = "tc"
			args = []string{"class", "del", "dev", intfName, "parent", htbRootHandle, "classid", htbRootClassid, "htb", "rate", rate, "ceil", rate}
			exeCMD(cmd, args)
		*/

	case "change":
		// change the br-int
		rate := nodeInbandwidth + "mbit"
		cmd := "tc"
		args := []string{"class", "change", "dev", brName, "parent", htbRootHandle, "classid", htbRootClassid, "htb", "rate", rate, "ceil", rate}
		exeCMD(cmd, args)

		// change the VM nodeDev
		rate = nodeOutbandwidth + "mbit"
		cmd = "tc"
		args = []string{"class", "change", "dev", intfName, "parent", htbRootHandle, "classid", htbRootClassid, "htb", "rate", rate, "ceil", rate}
		exeCMD(cmd, args)

	case "":
		log.Println("No change for ALL ", ip)

	default:

	}

	//configure default class
	ip = "default"
	action = ""
	rate := podDefaultInboundMin + "mbit"
	ceil := podDefaultInboundMax + "mbit"

	if _, ok := podQos[ip]; ok {
		rate = podQos[ip].InBandWidthMin + "mbit"
		ceil = podQos[ip].InBandWidthMax + "mbit"
		action = podQos[ip].Action
	}

	/*
		htbDefaultClassid = "8001"
		htbRootHandle = "1:"
		htbRootClassid = "1:1"
	*/
	htbDefaultClassidFull := htbRootHandle + htbDefaultClassid

	switch action {

	case "add":
		// br-int
		cmd := "tc"
		args := []string{"class", "add", "dev", brName, "parent", htbRootClassid, "classid", htbDefaultClassid, "htb", "rate", rate, "ceil", ceil}
		exeCMD(cmd, args)
		// VM nodeDev
		log.Println(" Set VM interface Default")
		cmd = "tc"
		args = []string{"class", "add", "dev", intfName, "parent", htbRootClassid, "classid", htbDefaultClassidFull, "htb", "rate", rate, "ceil", ceil}
		exeCMD(cmd, args)

	case "delete":
		// br-int
		cmd := "tc"
		args := []string{"class", "del", "dev", brName, "parent", htbRootClassid, "classid", htbDefaultClassid, "htb", "rate", rate, "ceil", ceil}
		exeCMD(cmd, args)
		// vm intface
		cmd = "tc"
		args = []string{"class", "del", "dev", intfName, "parent", htbRootClassid, "classid", htbDefaultClassidFull, "htb", "rate", rate, "ceil", ceil}
		exeCMD(cmd, args)

	case "change":

		cmd := "tc"
		args := []string{"class", "change", "dev", brName, "parent", htbRootClassid, "classid", htbDefaultClassid, "htb", "rate", rate, "ceil", ceil}
		exeCMD(cmd, args)

	case "":
		log.Println("No change for DEFAULT", ip)

	default:

	}

	//set tc class and filter for each pod
	setBrVmClassFilter(brName, podQos, podInfoMap)

}

func setBrVmClassFilter(brName string, podQos map[string]qosPara,
	podInfoMap map[string]podMetadata) {

	intfName := nodeDev

	for ip, val := range podQos {

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
			//configure tc qdisc htb class for each pod on bridgeDev
			classid := getCassid(classidPool)
			classidPool = decClassidPool(classidPool)
			curClassid := htbRootHandle + strconv.Itoa(classid)

			//log.Println("Add class and filters with current classID: " + curClassid)

			log.Println("add class on dev: " + brName)
			//cmd := "tc"
			//args := []string{"class", "add", "dev", brName, "parent", htbRootClassid, "classid",
			//	curClassid, "htb", "rate", rate, "ceil", ceil, "prio", prio}
			//exeCMD(cmd, args)
			cmd := " tc class add dev " + brName + " parent " + htbRootClassid + " classid " + curClassid + " htb rate " + rate + " ceil " + ceil + " prio " + prio
			exeCMDFull(cmd)

			//set tc filter for each pod on bridgeDev
			//tc filter add dev $nic parent 1:1 protocol ip prio 0 u32 match ip dst 10.0.3.153/32 flowid 1:2
			//log.Println(classid, curClassid)
			log.Println("add filter on dev: " + brName)
			//cmd = "tc"
			//args = []string{"filter", "add", "dev", brName, "parent", htbRootClassid, "protocol", "ip",
			//	"prio", "0", "u32", "match", "ip", "dst", ip + "/32", "flowid", curClassid}
			//exeCMD(cmd, args)
			cmd = " tc filter add dev " + brName + " parent " + htbRootClassid + " protocol ip prio 0 u32 match ip dst " + ip + "/32" + " flowid " + curClassid
			exeCMDFull(cmd)

			//get filter pref,
			//cmd = "tc"
			//args = []string{"filter", "show", "dev", brName, "parent", htbRootClassid}
			//output := exeCMD(cmd, args)

			cmd = " tc filter show dev " + brName + " parent " + htbRootClassid
			output := exeCMDFull(cmd)

			var pref string
			for _, line := range strings.Split(output, "\n") {

				if strings.Contains(line, curClassid) {
					pref = strings.Split(line, " ")[4]
				}
			}

			// config VM interface
			log.Println("add class on dev: " + intfName)
			//cmd = "tc"
			//args = []string{"class", "add", "dev", intfName, "parent", htbRootClassid, "classid",
			//	curClassid, "htb", "rate", rate, "ceil", ceil, "prio", prio}
			//exeCMD(cmd, args)

			cmd = " tc class add dev " + intfName + " parent " + htbRootClassid + " classid " + curClassid + " htb rate " + rate + " ceil " + ceil + " prio " + prio
			exeCMDFull(cmd)

			//filter cmd is "sudo tc filter add dev ens3 parent 1:0 bpf bytecode \"11,40 0 0 12,21 0 8 2048,48 0 0 23,21 0 6 17,40 0 0 42,69 1 0 2048,6 0 0 0,32 0 0 76,21 0 1 167838213,6 0 0 262144,6 0 0 0,\" flowid 1:100"

			//This is no VNI check, but check UDP port and src
			// 10,40 0 0 12,21 0 7 2048,48 0 0 23,21 0 5 17,40 0 0 36,21 0 3 4789,32 0 0 76,21 0 1 168431877,6 0 0 262144,6 0 0 0,

			// get the byte code
			bytecode := genBpfCode(ip)
			// filter's prio is different from class filter.
			// prio of filter means the seqence to check the filter. set 0 here and system will allocate the preference automatically
			// Take care, here I use 'htp_root_handle' (1:0) while not 'htbRootClassid'(1:1), or the traffic cannot match the filter

			log.Println("add filter on dev: " + intfName)
			filterCmd := "tc filter add dev " + string(intfName) + " parent " + htbRootHandle + " prio " + pref + " bpf bytecode " + bytecode + " flowid " + curClassid
			exeCMDFull(filterCmd)

			//update classid in podInfoMap
			if _, ok := podInfoMap[ip]; ok {
				podMeta := podInfoMap[ip]
				podMeta.classid = classid
				podMeta.pref = pref
				delete(podInfoMap, ip)
				podInfoMap[ip] = podMeta

			} else {
				log.Println("can not find in pod info map.", ip)
			}

		case "delete":
			if _, ok := podInfoMap[ip]; ok {
				classid := podInfoMap[ip].classid
				curClassid := htbRootHandle + strconv.Itoa(classid)
				pref := podInfoMap[ip].pref

				//log.Println("Delete pod filter on",curClassid, brName, ip, pref)

				cmd := "tc"
				args := []string{"filter", "del", "dev", brName, "parent", htbRootClassid, "prio", pref, "u32"}
				exeCMD(cmd, args)

				//log.Println("Delete pod class on",curClassid, brName, ip, pref)
				cmd = "tc"
				args = []string{"class", "del", "dev", brName, "parent", htbRootClassid, "classid",
					curClassid, "htb", "rate", rate, "ceil", ceil, "prio", prio}
				exeCMD(cmd, args)

				/*
					delete VM interface
					sudo tc filter del dev ens3 parent 1: prio 1
					sudo tc class del dev ens3 parent 1:1 classid 1:95
					note: there is a different with br-int that the prio and pref
				*/
				log.Println("delete filter on ", intfName)

				// Take care of 'htbRootHandle', not 'htbRootClassid' here
				cmd = "tc"
				args = []string{"filter", "del", "dev", intfName, "parent", htbRootHandle, "prio", pref}
				exeCMD(cmd, args)

				log.Println("delete class on ", intfName)
				cmd = "tc"
				args = []string{"class", "del", "dev", intfName, "parent", htbRootClassid, "classid",
					curClassid}
				exeCMD(cmd, args)

				// update classidPool
				classidPool = freeClassid(classid)

				podMeta := podInfoMap[ip]
				podMeta.classid = 0
				podMeta.pref = "0"
				podInfoMap[ip] = podMeta
				//fmt.Print(podInfoMap[ip])
			} else {
				log.Println("No config for " + ip + " in pod info map and ignore the DELETE action")
			}

		case "change":
			// There is 'change' cmd for linux TC and it takes affect when chaning the rate. However, it cannot take affect when we change one class's priority, unless we restart the traffic.
			// Here I take a temporary solution by firstly deleting the filter and class, then add again. Different with the 'delete' and 'add' cmd above, there is no need to update the classID from the class pool.
			log.Println("Before change, the podInfoMap is : ", podInfoMap)
			if _, ok := podInfoMap[ip]; ok {
				classid := podInfoMap[ip].classid
				curClassid := htbRootHandle + strconv.Itoa(classid)
				pref := podInfoMap[ip].pref
				log.Println("change class on" + brName + " and current classID is: " + curClassid)
				/*
					config br-int
				*/
				cmd := "tc"
				// first to delete the filter
				args := []string{"filter", "del", "dev", brName, "parent", htbRootClassid, "prio", pref, "u32"}
				exeCMD(cmd, args)

				// then delete the class
				args = []string{"class", "del", "dev", brName, "parent", htbRootClassid, "classid",
					curClassid, "htb", "rate", rate, "ceil", ceil, "prio", prio}
				exeCMD(cmd, args)

				// then add the new class (change)

				args = []string{"class", "add", "dev", brName, "parent", htbRootClassid, "classid",
					curClassid, "htb", "rate", rate, "ceil", ceil, "prio", prio}
				exeCMD(cmd, args)
				// finally add the filters

				args = []string{"filter", "add", "dev", brName, "parent", htbRootClassid, "protocol", "ip",
					"prio", pref, "u32", "match", "ip", "dst", ip + "/32", "flowid", curClassid}
				exeCMD(cmd, args)

				/*
					config vm
				*/
				// take care of parent, use 'htbRootHandle' while not 'htbRootClassid'
				cmd = "tc"
				// first del the filuer and class
				args = []string{"filter", "del", "dev", intfName, "parent", htbRootHandle, "prio", pref}
				exeCMD(cmd, args)

				args = []string{"class", "del", "dev", intfName, "parent", htbRootClassid, "classid",
					curClassid}
				exeCMD(cmd, args)

				// then add the new class and filter
				args = []string{"class", "add", "dev", intfName, "parent", htbRootClassid, "classid",
					curClassid, "htb", "rate", rate, "ceil", ceil, "prio", prio}
				exeCMD(cmd, args)
				// get the byte code
				bytecode := genBpfCode(ip)
				filterCmd := "tc filter add dev " + string(intfName) + " parent " + htbRootHandle + " prio " + pref + " bpf bytecode " + bytecode + " flowid " + curClassid
				exeCMDFull(filterCmd)

			} else {
				log.Println("Don't have this item, No change!")
			}

		case "":
			log.Println("No change of class and filters on br-int and vm intf on node: ", ip)

		default:

		}

		if len(classidPool) == 0 {
			log.Println("Error classid pool is empty. Cannot set pod bridge inbound bandwidth class and filter.")
			break
		}
	}

	//return classidPool
}

func genBpfCode(ip string) string {
	part1 := "\"10,40 0 0 12,21 0 7 2048,48 0 0 23,21 0 5 17,40 0 0 36,21 0 3 4789,32 0 0 76,21 0 1 "
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

func setVethTC(podQos map[string]qosPara, podInfoMap map[string]podMetadata) {
	//config pod inbound bandwidth tc qdisc on veth outside
	vethList := getVethList()
	for ip, val := range podQos {
		//log.Println(ip+" inbound: "+val.InBandWidthMin+", "+val.InBandWidthMax)
		//if ip = "all", skip it
		if ip == "all" || ip == "default" {
			//log.Println("all ip traffic bandwidth.")
			continue
		}

		//get container pid
		if podMeta, ok := podInfoMap[ip]; ok {
			containerPID := podMeta.cinfoList[0].pid
			action := val.Action

			switch action {
			case "add":
				//get veth id
				cmd := "nsenter"
				args := []string{"-t", containerPID, "-n", "ip", "link", "show", "eth0"}
				output := exeCMD(cmd, args)

				devID := strings.Split(output, ":")[0]
				tmp, _ := strconv.Atoi(devID)
				vethID := tmp + 1
				//log.Println("devID:",tmp, ", veth id:",strconv.Itoa(vethID))

				vethName := vethList[vethID]
				//log.Println(strconv.Itoa(vethID), vethName)

				//configure tc qdisc tbf on eth0 in the container(pod)
				//Firstly, delete the tc qdisc on the veth, tc qdisc del dev vethID root
				cmd = "tc"
				args = []string{"qdisc", "del", "dev", vethName, "root"}
				exeCMD(cmd, args)

				//secondly, configure tc qdisc tbf
				//tc qdisc add dev eth0 root tbf rate mbit latency 50ms burst 100k
				cmd = "tc"
				args = []string{"qdisc", "add", "dev", vethName, "root", "tbf", "rate",
					val.InBandWidthMax + "mbit", "latency", "50ms", "burst", "100k"}
				exeCMD(cmd, args)

				//show tc configuration
				log.Println("add qos on ", ip, vethName)
				//showTCQdisc(vethName)

			case "delete":
				//get veth id
				cmd := "nsenter"
				args := []string{"-t", containerPID, "-n", "ip", "link", "show", "eth0"}
				output := exeCMD(cmd, args)

				devID := strings.Split(output, ":")[0]
				tmp, _ := strconv.Atoi(devID)
				vethID := tmp + 1
				//log.Println("devID:",tmp, ", veth id:",strconv.Itoa(vethID))

				vethName := vethList[vethID]
				//log.Println(strconv.Itoa(vethID), vethName)

				//configure tc qdisc tbf on eth0 in the container(pod)
				//Firstly, delete the tc qdisc on the veth, tc qdisc del dev vethID root
				cmd = "tc"
				args = []string{"qdisc", "del", "dev", vethName, "root"}
				exeCMD(cmd, args)

				//show tc configuration
				log.Println("delete qos on ", ip, vethName)
				//showTCQdisc(vethName)

			case "change":
				//get veth id
				cmd := "nsenter"
				args := []string{"-t", containerPID, "-n", "ip", "link", "show", "eth0"}
				output := exeCMD(cmd, args)

				devID := strings.Split(output, ":")[0]
				tmp, _ := strconv.Atoi(devID)
				vethID := tmp + 1
				//log.Println("devID:",tmp, ", veth id:",strconv.Itoa(vethID))

				vethName := vethList[vethID]
				log.Println("Change pod Qos on", ip, vethName)
				cmd = "tc"
				args = []string{"qdisc", "change", "dev", vethName, "root", "tbf", "rate",
					val.InBandWidthMax + "mbit", "latency", "50ms", "burst", "100k"}
				exeCMD(cmd, args)

				//show tc configuration
				log.Println("change qos on ", ip, vethName)
				//showTCQdisc(vethName)

			case "":
				log.Println("Not change Qos on pod", ip)

			default:

			}

		} else {
			log.Println("Can not find ip: " + ip + " in podInfoMap.")
			continue
		}
	}
}

func setPodTC(podQos map[string]qosPara, podInfoMap map[string]podMetadata) {
	//config pod outbound bandwidth tc qdisc on eth0 in pod
	//log.Println("Start to set pod outbound bandwidth in pod...")
	for ip, val := range podQos {
		log.Println(ip + " outbound: " + val.OutBandWidthMin + ", " + val.OutBandWidthMax)
		//if ip = "all", skip it
		if ip == "all" || ip == "default" {
			//log.Println("all ip traffic bandwidth.")
			continue
		}
		//get container pid
		if podMeta, ok := podInfoMap[ip]; ok {
			containerPID := podMeta.cinfoList[0].pid
			action := val.Action

			switch action {
			case "add":
				//configure tc qdisc tbf on eth0 in the container(pod)
				//Firstly, delete the tc qdisc on the eth0, tc qdisc del dev $1 root
				//cmd := "nsenter"
				//args := []string{"-t", containerPID, "-n", "tc", "qdisc", "del", "dev", podDev, "root"}
				//exeCMD(cmd, args)

				//secondly, configure tc qdisc tbf
				//tc qdisc add dev eth0 root tbf rate mbit latency 50ms burst 100k
				log.Println("add tc for Pod:" + ip)
				cmd := "nsenter"
				args := []string{"-t", containerPID, "-n", "tc", "qdisc", "add", "dev", "eth0", "root", "tbf", "rate",
					val.OutBandWidthMax + "mbit", "latency", "50ms", "burst", "100k"}
				exeCMD(cmd, args)
				log.Println("Add qos on pod", ip, podDev)
				//showTCQdiscInPod(containerPID, podDev)

			case "delete":
				log.Println("Delete pod Qos on", ip, podDev)
				cmd := "nsenter"
				args := []string{"-t", containerPID, "-n", "tc", "qdisc", "del", "dev", podDev, "root"}
				exeCMD(cmd, args)
				//showTCQdiscInPod(containerPID, podDev)

			case "change":
				log.Println("Change pod Qos on", ip, podDev)

				cmd := "nsenter"
				args := []string{"-t", containerPID, "-n", "tc", "qdisc", "change", "dev", "eth0", "root", "tbf", "rate",
					val.OutBandWidthMax + "mbit", "latency", "50ms", "burst", "100k"}
				exeCMD(cmd, args)
				//showTCQdiscInPod(containerPID, podDev)

			case "":
				log.Println("Not change Qos on pod", ip, podDev)

			default:

			}
			//show tc configuration
			//log.Println("pod "+ip+" tc qdisc show: ")
			//showTCQdiscInPod(containerPID, podDev)
		} else {
			log.Println("Can not find ip: " + ip + " in podInfoMap.")
			continue
		}
	}
}

//get veth list on the host
func getVethList() map[int]string {
	//log.Println("Start to get veth list...")
	result := map[int]string{}
	intfList, err := net.Interfaces()
	if err != nil {
		log.Println(err)
	}
	for _, f := range intfList {
		// for flannel, it uses 'veth'
		//veth_key := "_l"
		//if strings.Contains(f.Name, veth_key) {
		result[f.Index] = f.Name
		//log.Println(f.Index, f.Name)
		//}
	}
	return result
}

//get the bridge int
func getBridgeName() string {
	bridgeDevf := "br"
	topoIntfNameMap := map[string]string{}
	topoIntfNameMap["br"] = "br"
	topoIntfNameMap["overlay_bridgeDev"] = "overlay_bridgeDev"

	ifaces, err := net.Interfaces()

	if err != nil {
		log.Printf("Error when decode interface %s\n", err)
	}

	for _, i := range ifaces {
		log.Printf(i.Name)
		_, ok := topoIntfNameMap[i.Name]
		if ok {
			//log.Printf("finally, we get %s\n", i.Name)
			return bridgeDevf
		}
	}
	return bridgeDevf
}

//get interface and it's IP address on VM
func getIntfIPaddress(intfName string) net.IP {
	var result net.IP
	result = nil
	ifaces, err := net.Interfaces()

	if err != nil {
		log.Printf("Error when decode interface %s\n", err)
	}
	// handle err
	for _, i := range ifaces {
		if i.Name == intfName {
			addrs, err := i.Addrs()
			// handle err
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
					log.Println("IP net is: ", ip)
					return ip
				case *net.IPAddr:
					result = v.IP
					ip = v.IP
					log.Println("IP address is: ", ip)
				}
				// process IP address
			}

			if err != nil {
				log.Printf("Error when decode interface %s\n", err)
			}
		}
	}
	return result
}
func initClassidPool(classidMax int) []int {
	var poolTemp []int
	for i := 2; i <= classidMax; i++ {
		poolTemp = append(poolTemp, i)
	}
	classidPool = poolTemp
	return classidPool
}

func getCassid(classidPool []int) int {

	len := len(classidPool)

	if len >= 3 {
		classid := classidPool[len-1]
		return classid
	}
	return 0
}

func decClassidPool(classidPool []int) []int {

	len := len(classidPool)

	if len >= 3 {
		return classidPool[:len-1]
	}
	return classidPool
}

func freeClassid(classid int) []int {

	classidPool = append(classidPool, classid)
	return classidPool
}

func exeCMDFull(cmd string) string {
	log.Println("full command is : ", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	//_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Println("Error to exec full CMD", cmd)
	}
	//log.Println("Output of command:", string(out))
	return string(out)
}

func exeCMD(cmd string, args []string) string {
	log.Println("command is ", cmd, " ", args)
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Printf("Error to exec CMD", cmd)
	}
	s := string(out)
	return s
}

func showTCQdisc(devName string) {

	//show tc configuration
	cmd := "tc"
	args := []string{"qdisc", "show", "dev", devName}
	log.Println(exeCMD(cmd, args))

}

func showTCClass(devName string) {

	//show tc configuration
	cmd := "tc"
	args := []string{"class", "show", "dev", devName}
	log.Println(exeCMD(cmd, args))

}

func showTCQdiscStatistics(devName string) {

	cmd := "tc"
	args := []string{"-s", "qdisc", "show", "dev", devName}
	log.Println(exeCMD(cmd, args))

}

func showTCClassStatistics(devName string) {

	cmd := "tc"
	args := []string{"-s", "class", "show", "dev", devName}
	log.Println(exeCMD(cmd, args))

}

func showTCQdiscInPod(containerPID string, devName string) {

	//show tc configuration
	cmd := "nsenter"
	args := []string{"-t", containerPID, "-n", "tc", "qdisc", "show", "dev", devName}
	log.Println(exeCMD(cmd, args))

}

func showTCClassInPod(containerPID string, devName string) {

	//show tc configuration
	cmd := "nsenter"
	args := []string{"-t", containerPID, "-n", "tc", "class", "show", "dev", devName}
	log.Println(exeCMD(cmd, args))

}

func showTCQdiscStatisticsInPod(containerPID string, devName string) {

	//show tc configuration
	cmd := "nsenter"
	args := []string{"-t", containerPID, "-n", "tc", "-s", "qdisc", "show", "dev", devName}
	log.Println(exeCMD(cmd, args))

}

func showTCClassStatisticsInPod(containerPID string, devName string) {

	//show tc configuration
	cmd := "nsenter"
	args := []string{"-t", containerPID, "-n", "tc", "-s", "class", "show", "dev", devName}
	log.Println(exeCMD(cmd, args))

}

func showTCFilter(devName string, handle string) {

	//show tc configuration
	cmd := "tc"
	args := []string{"filter", "show", "dev", devName, "parent", handle}
	log.Println(exeCMD(cmd, args))

}
