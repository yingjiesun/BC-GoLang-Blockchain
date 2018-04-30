package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"github.com/davecgh/go-spew/spew"
	//"github.com/joho/godotenv"
)

var Blockchain []Block
var temp_trans []Transaction
var	ext_ip = GetOutboundIP()
var	lan_ip = GetLAN_IP()
var brd_type = BROADCAST_IP_TYPE
var production_ip = ""

/*
YS: New array: ip_pool_dynamic, IPs of all nodes.
Every 30 minutes, external IP shall be sent to all nodes. The received IP shall be added to this array.
New joined node shall use ip_pool_dynamic first, if ip_pool_dynamic is empty, use the hard coded IP_POOL
If IP from IP_POOL is not reachable, enter IP manually (need new function)
*/

var ip_pool_dynamic []string

var t = time.Now()
var genesisBlock_data = []Transaction {	Transaction{ TransactionId: "This is Genesis Blok!"	, Timestamp: t.String()}}

func main() {
/*
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	*/

	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case sig := <-c:
			fmt.Printf("Got %s signal. Saving %d blocks to blockchain.dat \n", sig, len(Blockchain))
			blockChainPersisten("blockchain.dat")
			os.Exit(1)
		}
	}()

	if (brd_type == "LAN"){
		production_ip = lan_ip
	} else {
		production_ip = ext_ip
	}
	go func() {
		t := time.Now()
		genesisBlock := Block{0, production_ip, t.String(), genesisBlock_data, "", "", 0}
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
		Runtcp()
		//broadcast_IP(ext_ip)
	}()
	log.Fatal(run())
}
