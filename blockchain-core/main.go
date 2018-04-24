package main

import (
	
	"log"
	"time"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

var Blockchain []Block
var t = time.Now()
var genesisBlock_data = []Transaction {	Transaction{ TransactionId: "This is Genesis Blok!"	, Timestamp: t.String()} }

//YS: to hold transactions that are not saved in block yet
var temp_trans []Transaction 

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	
	go func() {
		t := time.Now()
		genesisBlock := Block{0, t.String(), genesisBlock_data, "", "", 0}
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)		
		
		Runtcp();
		
	}()
	log.Fatal(run())
	
}
