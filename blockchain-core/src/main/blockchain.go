package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"time"
	//	"os"
	"fmt"
//	"strconv"
	"sync"
	"errors"
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
var longest_chain []Block
var nounce int // in calculateHash() reset to -1 and calculate
var mutex = &sync.Mutex{}
var quit_hash_calc = make(chan bool) // quit_hash_calc <- true to terminate hash calculation

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

// check the whole chain and make sure it is self-consistent
func isChainValid(aChain []Block) bool {
	for i:= range aChain {
		if (i == 0) {
			continue
		}
		if !(isBlockValid(aChain[i], aChain[i - 1])) {
			return false
		}
	}
	return true
}

// make sure the chain we're checking is longer than the current blockchain
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		mutex.Lock()
		Blockchain = newBlocks
		mutex.Unlock()
		fmt.Println("Blockchain replaced by longer one")
	}
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
		select{
		case <- quit_hash_calc:
			fmt.Println("*** quit_hash_calc: ", block.Node_ip, " #", block.Index)
			return "Terminated"
		default:
			nounce++
			record := string(block.Index) + block.Timestamp + string(block_data) + block.PrevHash + string(nounce)
			h := sha256.New()
			h.Write([]byte(record))
			hashed := h.Sum(nil)
			returnValue = hex.EncodeToString(hashed)
			currentLeading = string(returnValue[0:difficulty])
		}
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
	t := time.Now().UTC()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Node_ip = node_ip
	newBlock.Timestamp = t.String()
	newBlock.Transactions = transactions
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	newBlock.Nounce = nounce
	if ( newBlock.Hash == "Terminated"){
		fmt.Println("generateBlock TERMINATED: ", newBlock.Node_ip, " #", newBlock.Index)
		error := errors.New("Terminated")
		return newBlock, error
	} else{
		return newBlock, nil
	}
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

//YS: only one combine of chain possible, since this function is called every time
func combine_chains()  {
	is_done := false
	if (len(potential_chains) > 1){
		for i := range potential_chains{
			if (is_done){
				break
			}
			for j := range potential_chains{
				if (potential_chains[i][len(potential_chains[i])-1].Hash == potential_chains[j][0].PrevHash ){

					potential_chains[i] = append(potential_chains[i], potential_chains[j]...)
					potential_chains = append(potential_chains[:j],potential_chains[j+1:]...)

					is_done = true
					break
				}
			}
		}
	}
}

func discard_short_chains(){
	if len(longest_chain) > 1 {
		highest_index := longest_chain[len(longest_chain) - 1].Index
		deleted := 0
		for k := range potential_chains{
			n := k - deleted
			if ((highest_index - potential_chains[n][len(potential_chains[n])-1].Index) > MAXFORKLENGTH ){
			//if (longest - len(potential_chains[n]) > 6){
				fmt.Println("Discard chain: ", n)
				potential_chains = append(potential_chains[:n], potential_chains[n+1:]...)
				deleted++
			}
		}
	}
}

// goLang slice is reference, this func is to create new array
func make_chain_from_slice(from_slice []Block) []Block{
	var new_chain []Block
	for i := range from_slice {
		new_chain = append(new_chain, from_slice[i])
	}
	return new_chain
}


/*
YS: Consensus
When a new block received, it could be
a. A child of a chain's youngest block - attach
b. A child of a chain's middle/oldest block - create new chain with the section of parent and older blocks, attach
c. A parent of a chain's oldest block - attach
d. A parent of a chain's middle/youngest block? - NOT POSSIBLE, means too blocks have same hash
e. none of above - block is orphan, or the very first block, create new chain with this orphan block
f. combination of a, b, c - After the above actions, try to connect chains and discard old ones
g. trim chains, discard short chains
h. handle confirmed section of blockchain, save to hard drive
*/

func add_block_to_potential_chains(){
		for {
			select{
				case bl, ok := <- waiting_blocks:
					if ok{
							connected := false
							mutex.Lock()
							if (len(potential_chains) == 0){
								fmt.Println("This is the very first block" )
								potential_chains = [][]Block{{bl}}
								connected = true
							//	fmt.Println("potential_chains Length: ", len(potential_chains) )
							} else if ( (longest_chain[len(longest_chain) - 1].Index - bl.Index) <= MAXFORKLENGTH ) {
								for i := range potential_chains{
									if (connected) {
										break
									}
									for j := range potential_chains[i] {
										if (connected) {
											break
										}
										if (potential_chains[i][j].Hash == bl.PrevHash){ //parent found
											if (j == (len(potential_chains[i])-1)){ //parent is last block
												fmt.Println("Added block as child: ", bl.Node_ip, "  #", bl.Index )
												potential_chains[i] = append(potential_chains[i], bl)
												connected = true
												break
											} else { //parent is in middle of chain
												fmt.Println("Created new chain, block added as child: ",  bl.Node_ip  , " #",bl.Index)
												new_chain := make_chain_from_slice(potential_chains[i][:j+1]) //j+1, so that j is included in new_chain
												new_chain = append(new_chain,bl)
												potential_chains = append(potential_chains, new_chain)
												connected = true
												break
											}
										}
									}
								}

								if (!connected) { // no parent found, finding child
									for k := range potential_chains{
										if (potential_chains[k][0].PrevHash == bl.Hash){ //child found
											potential_chains[k] = append([]Block{bl},potential_chains[k]...)
											fmt.Println("Added block as PARENT, ", bl.Node_ip , " #", bl.Index)
											connected = true
											break
										}
									}
								}
								if (!connected){
									fmt.Println("Add block as ORPHAN:", bl.Node_ip , " #", bl.Index  )
									potential_chains = append(potential_chains, []Block{bl})
								}
							}
							combine_chains()
							discard_short_chains()
							setLongestChain()
							mutex.Unlock()
							for k := range potential_chains{
									fmt.Println("chain #", k," : ", len(potential_chains[k]) )
							}
							is_new_block_added = true
						} else {
							//	fmt.Println("waiting_blacks channel closed")
						}
			default:
			//	fmt.Println("No waiting block, moving on")
		}
	}
}

func setLongestChain() {
		longest_index := 0
		for i := range potential_chains{
			if (len(potential_chains[i]) > len(potential_chains[longest_index])){
				longest_index = i
			}
		}
		fmt.Println("Longest chain is: " , longest_index )
		longest_chain = potential_chains[longest_index]
}
