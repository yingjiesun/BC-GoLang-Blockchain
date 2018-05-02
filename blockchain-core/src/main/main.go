package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"github.com/davecgh/go-spew/spew"
//	"strconv"
	//"github.com/joho/godotenv"
)

var Blockchain []Block
var temp_trans []Transaction
var	ext_ip = GetOutboundIP()
var	lan_ip = GetLAN_IP()
var brd_type = BROADCAST_IP_TYPE
var production_ip = ""

/*
YS: New array: peer_ip_pool, IPs of all nodes.
Every 30 minutes, external IP shall be sent to all nodes. The received IP shall be added to this array.
New joined node shall use peer_ip_pool first, if peer_ip_pool is empty, use the hard coded IP_POOL
If IP from IP_POOL is not reachable, enter IP manually (need new function)
*/






var peer_ip_pool []string

var t = time.Now()
var genesisBlock_data = []Transaction {	Transaction{ TransactionId: "This is Genesis Blok!"	, Timestamp: t.String()}}

func main() {


	var t1 = time.Now().Unix()
//	var t1_str = strconv.Itoa(t)
	fmt.Println("Unix: " , t1)


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

	//YS: add its own ip to peer_ip_pool
	peer_ip_pool = append_if_missing(peer_ip_pool, production_ip)

	for v := range SEED_IP_POOL{
			peer_ip_pool = append_if_missing(peer_ip_pool, SEED_IP_POOL[v])
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
