package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
//	"github.com/davecgh/go-spew/spew"
//	"strconv"
	//"github.com/joho/godotenv"
)
var genesisBlock Block
var Blockchain []Block
var temp_trans []Transaction
var	ext_ip = GetOutboundIP()
var	lan_ip = GetLAN_IP()
var brd_type = BROADCAST_IP_TYPE
var production_ip = ""
var waiting_blocks = make(chan Block)
/*
YS: New array: peer_ip_pool, IPs of all nodes.
Every 30 minutes, external IP shall be sent to all nodes. The received IP shall be added to this array.
New joined node shall use peer_ip_pool first, if peer_ip_pool is empty, use the hard coded IP_POOL
If IP from IP_POOL is not reachable, enter IP manually (need new function)
*/
var peer_ip_pool []string

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
	//YS: add its own ip to peer_ip_pool
	peer_ip_pool = append_if_missing(peer_ip_pool, production_ip)
// add seed to perr IP pool
	for v := range SEED_IP_POOL{
			peer_ip_pool = append_if_missing(peer_ip_pool, SEED_IP_POOL[v])
	}
	fmt.Println("==========================")
	fmt.Println("My External IP: " + ext_ip );
	fmt.Println("My Internal IP: " + lan_ip );
	fmt.Println("My Production IP: " + production_ip );
	for v := range peer_ip_pool{
		fmt.Println("peer_ip_pool : " + peer_ip_pool[v])
	}
	t := time.Now().UTC()
	tran_id := "0, " + production_ip
	msg := "This is genesis block"
	tranaction_genesis := Transaction{ TransactionId: tran_id, Timestamp: t.String(), Msg: msg, TransactionHash: ""}
	tranaction_genesis.TransactionHash = get_tran_hash(tranaction_genesis)
	var genesisBlock_data = []Transaction {	tranaction_genesis }
	genesisBlock = Block{0, production_ip, t.String(), genesisBlock_data, "", "", 0}
	genesisBlock.Hash = calculateHash(genesisBlock)
	genesisBlock.Nounce = nounce
//	spew.Dump(genesisBlock)
	go block_to_chan(genesisBlock, waiting_blocks)
	go add_block_to_potential_chains()
	fmt.Println("after add")
	go Runtcp()
	fmt.Println("after Runtcp")
	log.Fatal(run())
}

func block_to_chan(blk Block, c chan Block) {
	c <- blk // send blk to c
}
