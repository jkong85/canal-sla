//bandwidth control interface for EyeQ
//author: Yan Sun

package main

import (
	//"os"
	"encoding/json"
	//"io/ioutil"
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"golang.org/x/net/context"
	"log"
	"net"
	"time"
)

type ContainerBW struct {
	NodeIP          string
	PodID           string
	VlanID          string
	VxlanID         string
	PodIP           string
	Action          string
	InBandWidthMin  string // unit is Mbps
	InBandWidthMax  string // unit is Mbps
	OutBandWidthMin string // unit is Mbps
	OutBandWidthMax string // unit is Mbps
	PodPriority     string // 0-7, 0 is the highest priority, 7 is the lowest priority.
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
func SaveAsJson(bw []ContainerBW, path string) {
	b, err := json.Marshal(bw)
	check(err)
	ioutil.WriteFile(path, b, 0644)
}
*/
func main() {

	l, err := net.Interfaces()
	if err != nil {
		println(err)

	}
	for _, f := range l {
		fmt.Println(f.Index, f.Name)
	}

	bw := []ContainerBW{
		{"192.0.0.1", "1", "100", "1", "all", "", "1000", "1000", "1000", "1000", "0"},
		{"192.0.0.1", "1", "100", "1", "default", "", "10", "100", "10", "100", "0"},
		{"192.0.0.1", "1", "100", "1", "10.1.2.5", "delete", "10", "60", "10", "60", "0"},
		{"192.0.0.1", "1", "100", "1", "10.1.2.6", "delete", "10", "60", "10", "60", "5"},
		{"192.0.0.1", "1", "100", "1", "10.1.2.7", "delete", "10", "60", "10", "60", "7"},
	}

	// We can register as many backends that are supported by libkv
	etcd.Register()

	server := "127.0.0.1:4001"

	// Initialize a new store with consul
	kv, err := libkv.NewStore(
		store.ETCD,
		[]string{server},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
	if err != nil {
		log.Fatal("Cannot create store", kv)
	}

	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:4001"},
		Transport: client.DefaultTransport,
		//set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := client.NewKeysAPI(c)

	qos_encode, err := json.Marshal(bw)
	if err != nil {
		log.Fatal(err)
	}

	intf, err := net.InterfaceByName("ens3")

	if err != nil {
		log.Fatal("Cannot find interface by name eth0")
	}

	mac := intf.HardwareAddr

	key := "/" + string(mac)

	resp, err := kapi.Set(context.Background(), key, string(qos_encode), nil)
	if err != nil {
		log.Fatal(err)
	} else {
		//log.Printf("Set is done. Metadata is %q\n", resp)
		fmt.Println("Set is done. Metadata is %q\n", resp)
	}

	/*
		path := "./qos.json"
		f, err := os.Create(path)
		check(err)

		defer f.Close()
		SaveAsJson(bw, path)
		f.Sync()
	*/
}
