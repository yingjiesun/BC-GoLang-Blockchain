package main

type Block struct {
	Index     int
	Timestamp string
	Transactions       []Transaction
	Hash      string
	PrevHash  string
	Nounce    int
}