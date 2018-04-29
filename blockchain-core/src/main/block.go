package main

//YS 4/29 added node_ip
type Block struct {
	Index     int
	Node_ip   string
	Timestamp string
	Transactions       []Transaction
	Hash      string
	PrevHash  string
	Nounce    int
}
