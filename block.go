package main

type Block struct {
	Index     int
	Timestamp string
	//BPM       int
	BPM       []Transaction
	Hash      string
	PrevHash  string
}