package main

import (
	
	"crypto/sha256"
	"encoding/hex"
	"time"
	"encoding/json"
	"os"
	"strconv"
)

//YS: nounce will get value in calculateHash()
var nounce int

// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// make sure the chain we're checking is longer than the current blockchain
func replaceChain(newBlocks []Block) {
	mutex.Lock()
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
	mutex.Unlock()
}

/*

YS: 
Add Nounce to generate qualified hash, DIFFICULTY is from .env file
(GoLang allows a function to return two values, Nounce and hash can be returned at same time)

*/

// SHA256 hasing
func calculateHash(block Block) string {

	returnValue := "NOT ME"
	difficulty, err := strconv.Atoi(os.Getenv("DIFFICULTY"))
	nounce = -1
	
	a := &block.Transactions
	block_data, err := json.Marshal(a)
	if err != nil {
        panic (err)
    }
	
	requiredLeadings := getRequiredString(difficulty)
	currentLeading := "XXXXXXXXXXXXXXXXXXXXXXXXX"
	for currentLeading != requiredLeadings {
		nounce++
		record := string(block.Index) + block.Timestamp + string(block_data) + block.PrevHash + string(nounce)
		h := sha256.New()		
		h.Write([]byte(record))
		hashed := h.Sum(nil)
		returnValue = hex.EncodeToString(hashed)
		currentLeading = string(returnValue[0:difficulty])		
	}
	return returnValue
}

//YS: generate string of required leading 0s 
func getRequiredString(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = '0'
    }
    return string(b)
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, transactions []Transaction) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Transactions = transactions
	newBlock.PrevHash = oldBlock.Hash	
	newBlock.Hash = calculateHash(newBlock)	
	newBlock.Nounce = nounce

	return newBlock, nil
}