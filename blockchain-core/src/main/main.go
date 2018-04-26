package main

import (
	"log"
	"time"
	"github.com/davecgh/go-spew/spew"
	//"github.com/joho/godotenv"
)

var Blockchain []Block
var temp_trans []Transaction

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
	go func() {
		t := time.Now()
		genesisBlock := Block{0, t.String(), genesisBlock_data, "", "", 0}
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
		Runtcp();
	}()
	log.Fatal(run())
}
