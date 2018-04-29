/**
YS: Handle transaction in temp_trans
*/

package main

import (
	"time"
	"fmt"
)

//YS: this should trigger propagation of transactions

func append_temp_trans(new_tx Transaction){
	temp_trans = append(temp_trans, new_tx)
	//todo: should check duplicate in case
}

//YS: delete transactions that are older than the timestamp of last block
//YS: not functioning yet

func update_temp_trans(){
	last_block_time := Blockchain[len(Blockchain)-1].Timestamp
	block_time, err := time.Parse(time.RFC3339,last_block_time)
	if err != nil {
		panic(err)
		fmt.Println("ERROR update_temp_trans()")
	}
	for i := range temp_trans {
		tx_time, err := time.Parse(time.RFC3339, temp_trans[i].Timestamp)
		if err != nil {
			panic(err)
			fmt.Println("ERROR update_temp_trans()")
		}
		if (tx_time.Before(block_time)){
			//YS: this is wrong, array size changed
			temp_trans = remove_tx_by_index(temp_trans, i)
		}
	}
}

func remove_tx_by_index(s []Transaction, i int) []Transaction {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}
