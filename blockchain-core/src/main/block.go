package main

import (
	"encoding/gob"
	"os"
)

type Block struct {
	Index        int
	Node_ip      string
	Timestamp    string
	Transactions []Transaction
	Hash         string
	PrevHash     string
	Nounce       int
}

func (block *Block) persistent(file *os.File) {
	if file != nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(block)
	}
}
