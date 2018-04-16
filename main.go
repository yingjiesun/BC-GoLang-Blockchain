package main

import (
	
	"log"
	"time"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

var Blockchain []Block

var genesisBlock_data = []Transaction {	Transaction{ transactionId: "This is Genesis Blok!"	} }

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	
	go func() {
		t := time.Now()
		genesisBlock := Block{0, t.String(), genesisBlock_data, "", ""}
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)		
		
		Runtcp();
		
	}()
	log.Fatal(run())
	
}

