package main

import (
	
	"crypto/sha256"
	"encoding/hex"
	"time"
	"encoding/json"
)


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

// SHA256 hasing
func calculateHash(block Block) string {

	a := &block.BPM
	block_data, err := json.Marshal(a)
	if err != nil {
        panic (err)
    }

	record := string(block.Index) + block.Timestamp + string(block_data) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, BPM []Transaction) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}