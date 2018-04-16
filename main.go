package main

import (
	
	"log"
	"time"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

var Blockchain []Block

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	
	go func() {
		t := time.Now()
		genesisBlock := Block{0, t.String(), 0, "", ""}
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)		
		
		Runtcp();
		
	}()
	log.Fatal(run())
	
}

