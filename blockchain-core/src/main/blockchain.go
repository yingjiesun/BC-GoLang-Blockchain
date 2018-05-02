package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"time"
	//	"os"
	"fmt"
	"strconv"
	"sync"
)

/*
YS: Two dimensional array of blockchain, received blocks will be saved in potential_chains
When a block is received, first check if it can be connected to existing chain (parent or child),
if not, create new chain and put this orphan into it
calculate the maximum difference between chains, if the largest difference is greater than 6
dicard the shortest one, and save the oldese N blocks from the longest chaing
N is the length of the discarded chain
*/

var potential_chains [][]Block

//YS: nounce get value in calculateHash()
var nounce int
var mutex = &sync.Mutex{}

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

	fmt.Println("Current blockchain length: " + strconv.Itoa(len(Blockchain)))
	fmt.Println("New blockchain length: " + strconv.Itoa(len(newBlocks)))

	mutex.Lock()
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
		fmt.Println("Blockchain replaced by longer one")
	}
	mutex.Unlock()
}

// SHA256 hasing
func calculateHash(block Block) string {

	returnValue := "NOT ME"
	difficulty := DIFFICULTY
	nounce = -1

	a := &block.Transactions
	block_data, err := json.Marshal(a)
	if err != nil {
		panic(err)
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
func generateBlock(oldBlock Block, transactions []Transaction, node_ip string) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Node_ip = node_ip
	newBlock.Timestamp = t.String()
	newBlock.Transactions = transactions
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	newBlock.Nounce = nounce


	return newBlock, nil
}

func blockChainPersisten(path string) {
	file, err := os.Create(path)

	if err == nil {
		for i, block := range Blockchain {
			block.persistent(file)
			fmt.Print(i, block.Index)
		}
	} else {
		fmt.Printf("open file failed while saving\n")
	}
}

func add_block_to_potential_chains(bl Block){
	fmt.Println("=============START CONSENSUS===============" )

	connected := false;
	for i := range potential_chains{
		for j := range potential_chains[i] {
			if (potential_chains[i][j].Hash == bl.PrevHash ||
					potential_chains[i][j].PrevHash == bl.Hash){
						potential_chains[i] = append(potential_chains[i], bl)
						connected = true
						fmt.Println("Block is connected to a potential chain: ", i , ", Block from: ", bl.Node_ip )
						break
					}
		}
		fmt.Println("potential_chain ", i, " length: ", len(potential_chains[i]) )
	}

	if (!connected){
		new_chain := []Block{bl}
		potential_chains = append(potential_chains, new_chain)
		fmt.Println("Block is orphan, created new chain and appened to potential_chains, Block from: ", bl.Node_ip )
	}
	update_potential_chains()
}

func update_potential_chains(){
	fmt.Println("In function update_potential_chains()" )
	max_length := 10 //TODO: length of longest chain in update_potential_chains
	min_length :=0 //TODO: length of shortest chain in update_potential_chains
	if (max_length - min_length > MAXFORKLENGTH) {
		//TODO: save oldest min_length blocks of longest chain to hard drive, discard shortest chain
	}

		fmt.Println("=============END CONSENSUS===============" )
}
