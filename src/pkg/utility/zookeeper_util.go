package utility

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/UrlShortener/src/pkg/config"
	"github.com/go-zookeeper/zk"
)

const rangeBlock = 1000000

type SequenceRange struct {
	Start int
	Curr  int
	End   int
}

var (
	zkConnection *zk.Conn
	NodeRange    *SequenceRange
)

func GetHash(n int) string {
	hashRange := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	hashString := ""

	for n > 0 {
		hashString = string(hashRange[n%62]) + hashString
		n = n / 62
	}

	return hashString
}

func ConnectZookeeper() bool {
	servers := config.AppConfig.ZkServers
	conn, events, err := zk.Connect(servers, 300*time.Second)

	if err != nil {
		fmt.Println("Error while connecting to Zookeeper")
		return false
	}

	zkConnection = conn

	for e := range events {
		if e.State == zk.StateConnected {
			fmt.Printf("Connected %s\n", e.State)
			break
		}
	}

	fmt.Println(zkConnection.SessionID())

	return true
}

func CreateRangeNode(path string) {
	log.Printf("Create Path Request %s\n", path)
	if ExistsNode(path) {
		log.Printf("Path already exists %s\n", path)
		return
	}
	data := []byte(strconv.Itoa(0))
	acl := zk.WorldACL(zk.PermAll)
	s, err := zkConnection.Create(path, data, 0, acl)

	if err != nil {
		fmt.Println("Failed to create path")
		return
	}

	fmt.Printf("Path Created %s", s)

}

func GetRangeNode() {
	log.Printf("Get Next Range\n")
	data, stat, err := zkConnection.Get("/range")

	if err != nil {
		fmt.Printf("Error while getting range %s", err.Error())
		return
	}

	fmt.Println(stat.Version)
	d := string(data)
	startRange, err := strconv.Atoi(d)

	if err != nil {
		log.Panicf("Error while parsing %e", err)
		return
	}

	NodeRange.Start = startRange + rangeBlock
	NodeRange.Curr = startRange + rangeBlock
	NodeRange.End = NodeRange.Start + rangeBlock
	log.Printf("Received Range %s\n", d)
	setRangeNode(NodeRange.Start)
}

func setRangeNode(newRange int) {
	log.Printf("Setting New Range to %d\n", newRange)
	data := []byte(strconv.Itoa(newRange))
	stat, err := zkConnection.Set("/range", data, -1)

	if err != nil {
		log.Panicf("Failed to set data %s", err.Error())
		return
	}

	log.Printf("Range updated %d\n", stat.Version)
}

func ExistsNode(path string) bool {
	log.Printf("Checking Node %s Exists\n", path)
	exists, stat, err := zkConnection.Exists(path)

	if err != nil {
		log.Panicf("Failed to set data %s\n", err.Error())
		return false
	}

	log.Println(stat.Version)
	return exists
}

func RemoveNode(path string) {
	if err := zkConnection.Delete(path, -1); err != nil {
		fmt.Printf("Failed to delete node %s\n", path)
	}
}

func Close() {
	zkConnection.Close()
}
