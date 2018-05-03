/**
YS:
Create TCP server
Connect to IPs in IP pool, send blockchain
Listen port 9991 and accept all connection requests
Receive blockchain and call replaceChain function
*/

package main

import (
	//"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	//"os"
	"time"
	"fmt"
	"strings"
//	"github.com/davecgh/go-spew/spew"
	"strconv"
	"bytes"
//	"encoding/binary"
  "io/ioutil"
  "net/http"
	"math/rand"
)

// bcServer handles incoming concurrent Blocks
//var bcServer chan []Block

func Runtcp() error {

	server, err := net.Listen("tcp", ":" + TCP_PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	//YS: Print Node External and Internal IP
	fmt.Println("My External IP: " + ext_ip );
	fmt.Println("My Internal IP: " + lan_ip );

	/*
	YS: Testing code below, this is to generate block randomly in 30-90 seconds
	*/
	go func(){
		for {
			fmt.Println("Start creating new block")
			if len(temp_trans) > 0 {
				mutex.Lock()
				newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], temp_trans, production_ip )
			//	temp_trans = temp_trans[:0]
				temp_trans = nil
				mutex.Unlock()
				if err != nil {
					//panic (err)
					fmt.Println("Error creating new block")
				}
				if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
					newBlockchain := append(Blockchain, newBlock)
					replaceChain(newBlockchain)
					fmt.Println("====NEW BLOCK CREATED AND ADDED!====")
					//spew.Dump(Blockchain)

					add_block_to_potential_chains(newBlock)
					propagate_BL(newBlock)

				} else {
					fmt.Println("INVALID block")
				}
			}
			time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
		}
	}()

	//YS: END of generating test block

	/*
	YS: Testing code below, this is to generate transaction randomly in 10-20 seconds
	*/
	go func(){
		test_tran_id := 1
		for {
			t := time.Now()
			var tranaction_new = Transaction{ TransactionId: strconv.Itoa(test_tran_id) + ", created by " + production_ip, Timestamp: t.String()}
			mutex.Lock()
			append_temp_trans(tranaction_new)
			mutex.Unlock()
			test_tran_id++
			propagate_TX(tranaction_new)
			time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
		}
	}()

	//YS: END of generating test transaction

	go broadcast_IP(production_ip)
	go dialConn_bc()

	//YS: create a new connection each time a connection request is received, and serve it.
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
	return nil
}

//YS: Function to get the public ip address (IP of router)
//YS: currently this is only to display IP in termainl

func GetOutboundIP() string {
	rsp, err := http.Get("http://checkip.amazonaws.com")
    if err != nil {
        return "Get External IP Error 001"
    }
    defer rsp.Body.Close()
    buf, err := ioutil.ReadAll(rsp.Body)
    if err != nil {
        return "Get External IP Error 002"
    }
    return string(bytes.TrimSpace(buf))
}

//YS: Function to get IP of hosting device (IP of LAN 192.168...)

func GetLAN_IP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
				log.Fatal(err)
		}
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	fmt.Println("outer IP: " + localAddr[0:idx] );
    return localAddr[0:idx]
}

/*
YS 4/20:
Create TCP connection to other nodes listed in IP_POOL
Every 30 seconds, dial TCP and send blockchain
*/

func dialConn_bc() {

	//YS: this should be updated to use dynamic ip pool first
	//ip_pool := strings.Split( IP_POOL , "_")
	broadcast_invl := BROADCAST_INTERVAL

	for v := range peer_ip_pool{
		fmt.Println("peer_ip_pool : " + peer_ip_pool[v])
	}
	for {
		//YS: Prepare data
		raw_chain, err := json.Marshal(Blockchain)
		if ( err != nil ){
				fmt.Println("raw_chain Marshal error : " , err.Error())
		}
		json_data := Container{Type:"BC", Object:raw_chain}

		for i := range peer_ip_pool{
			//YS: do not dial itself
			if (peer_ip_pool[i] == production_ip) {
				continue
			}
			conn, err := net.Dial("tcp", peer_ip_pool[i] + ":" + TCP_PORT)
			if err != nil {
				fmt.Println("CONNECTION ERROR:", err.Error())
				continue
			}
			encoder := json.NewEncoder(conn)

			if err := encoder.Encode(json_data); err != nil {
					 fmt.Println("dialConn_bc() encode.Encode error: ", err)
			 }
		}
		time.Sleep(time.Duration(broadcast_invl) * time.Second)
	}
}


/*
YS: when node start, call this function once.
This function will send IP to known seed nodes in IP_POOL,
and when seed nodes receive this broadcast, the peer_ip_pool (peer nodes)
will be broadcast to all nodes
*/
func broadcast_IP(the_ip string) {

	raw_ip, err := json.Marshal(the_ip)
	if ( err != nil ){
			fmt.Println("raw_ip Marshal error : " , err.Error())
	}
	json_data := Container{Type:"IP", Object:raw_ip}

		for i := range SEED_IP_POOL{
			//YS: do not dial itself
			if (SEED_IP_POOL[i] == production_ip) {
				continue
			}
			conn, err := net.Dial("tcp", SEED_IP_POOL[i] + ":" + TCP_PORT)
			if err != nil {
				fmt.Println("CONNECTION ERROR:", err.Error())
				continue
			}

			encoder := json.NewEncoder(conn)
			if err := encoder.Encode(json_data); err != nil {
					 fmt.Println("broadcast_IP() encode.Encode error: ", err)
			 }
	}
}

func broadcast_IPS() {

	raw_ips, err := json.Marshal(peer_ip_pool)
	if ( err != nil ){
			fmt.Println("raw_ips Marshal error : " , err.Error())
	}
	json_data := Container{Type:"IPS", Object:raw_ips}

	for i := range peer_ip_pool{
			//YS: do not dial itself
			if (peer_ip_pool[i] == production_ip) {
				continue
			}
			conn, err := net.Dial("tcp", peer_ip_pool[i] + ":" + TCP_PORT)
			if err != nil {
				fmt.Println("CONNECTION ERROR:", err.Error())
				continue
			}

			encoder := json.NewEncoder(conn)
			if err := encoder.Encode(json_data); err != nil {
				 fmt.Println("broadcast_IPS() encode.Encode error: ", err)
		 }
	}
}

//YS: when new transaction generated, call this function.
func propagate_TX(the_tx Transaction) {

	//YS: this should be updated to use dynamic ip pool first
	raw_tx, err := json.Marshal(the_tx)
	if ( err != nil ){
			fmt.Println("raw_tx Marshal error : " , err.Error())
	}
	json_data := Container{Type:"TX", Object:raw_tx}

	for i := range peer_ip_pool{
		//YS: do not dial itself
		if (peer_ip_pool[i] == production_ip) {
			continue
		}
		conn, err := net.Dial("tcp", peer_ip_pool[i] + ":" + TCP_PORT)
		if err != nil {
			fmt.Println("CONNECTION ERROR:", err.Error())
			continue
		}
		encoder := json.NewEncoder(conn)
		if err := encoder.Encode(json_data); err != nil {
				 fmt.Println("broadcast_TX() encode.Encode error: ", err)
		}
	}
}

//YS: when new block generated, call this function.
func propagate_BL(the_bl Block) {

	raw_bl, err := json.Marshal(the_bl)
	if ( err != nil ){
			fmt.Println("raw_bl Marshal error : " , err.Error())
	}
	json_data := Container{Type:"BL", Object:raw_bl}

	for i := range peer_ip_pool{
		//YS: do not dial itself
		if (peer_ip_pool[i] == production_ip) {
			continue
		}
		conn, err := net.Dial("tcp", peer_ip_pool[i] + ":" + TCP_PORT)
		if err != nil {
			fmt.Println("CONNECTION ERROR:", err.Error())
			continue
		}
		encoder := json.NewEncoder(conn)
		if err := encoder.Encode(json_data); err != nil {
				 fmt.Println("broadcast_BL() encode.Encode error: ", err)
		}
	}
}

/*
YS: Accept TCP connection
*/
//YS: new function that accept different type of json objects, and decode accordingly

func handleConn(conn net.Conn) {

	var c Container
	defer conn.Close()

	decoder := json.NewDecoder(conn)

	if err := decoder.Decode(&c); ( err != nil && err != io.EOF ){
		fmt.Println("decoder.Decode(&c) err: ", err)
	} else {
		t := time.Now()
		fmt.Println("==========================")
		fmt.Println("My External IP: " + ext_ip );
		fmt.Println("My Internal IP: " + lan_ip );
		fmt.Println("Received data from other nodes: " + t.String())
		fmt.Println("c.Type: " + c.Type)

		switch c.Type {
		case "TX":
				process_TX(c)
		case "BC":
			process_BC(c)
		case "IP":
			process_IP(c)
		case "IPS":
			process_IPS(c)
		case "BL":
			process_BL(c)
		default:
			fmt.Println("Can not process data received")
		}
	//	spew.Dump(received_blockchain)
	}
}

func process_BC(c Container){
	var received_blockchain []Block
	json.Unmarshal(c.Object, &received_blockchain)
	replaceChain(received_blockchain)
//	fmt.Println("===========Received Blockchain:")
//	spew.Dump(received_blockchain)
//	fmt.Println("===========END Received Blockchain:")
}

/*
YS: Only seed nodes in IP_POOL receive this IP
IP will be added to peer_ip_pool and then peer_ip_pool will be broadcast to all nodes
*/
func process_IP(c Container){
	var received_ip string
	json.Unmarshal(c.Object, &received_ip)
	peer_ip_pool = append_if_missing(peer_ip_pool, received_ip)
	fmt.Println("Received IP added to dynamic ip pool: " + received_ip)
	for ele := range peer_ip_pool{
		fmt.Println("dynamic ip pool: " + peer_ip_pool[ele])
	}
	//YS: Seed node receive new node IP and broadcast to all nodes
	broadcast_IPS()
}

func process_IPS(c Container){
	var received_ips []string
	json.Unmarshal(c.Object, &received_ips)

	peer_ip_pool = merge_array_unique(peer_ip_pool, received_ips)

	fmt.Println("peer_ip_pool received and merged" )
	for ele := range peer_ip_pool{
		fmt.Println("dynamic ip pool: " + peer_ip_pool[ele])
	}
}

func process_TX(c Container){
	var received_tx Transaction
	json.Unmarshal(c.Object, &received_tx)
	temp_trans = append(temp_trans, received_tx)
	fmt.Println("Received TX added to temp_trans: " + received_tx.TransactionId)
}

func process_BL(c Container){
	var received_bl Block
	json.Unmarshal(c.Object, &received_bl)
	if (calculateHash(received_bl) == received_bl.Hash){
		add_block_to_potential_chains(received_bl)
	}

	fmt.Println("Received block" )
}

func append_if_missing(slice []string, str string) []string {
    for _, ele := range slice {
        if ele == str {
            return slice
        }
    }
    return append(slice, str)
}

func merge_array_unique(slice1 []string, new_slice []string) []string {
    for _, ele := range new_slice {
        slice1 = append_if_missing(slice1, ele)
    }
    return slice1
}

/*
//YS: Original hadleConn that take blockchain json obj only
func handleConn(conn net.Conn) {
	var received_blockchain []Block
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	decoder.Decode(&received_blockchain)
	t := time.Now()
	fmt.Println("===============================================")
	fmt.Println("My External IP: " + ext_ip );
	fmt.Println("My Internal IP: " + lan_ip );
	fmt.Println("Received blockchain: " + t.String())
	spew.Dump(received_blockchain)
	//YS: function in blockchain.go, replace with new chain if it is longer
	replaceChain(received_blockchain)
  conn.Close()
}
*/
/*
YS: Data_type is used for receiving node to identify what kind of data received
IP : Data_type = "IP"
transaction: Data_type = "TX"
Blockchain : Data_type = "BC"
*/

type Container struct {
    Type   string
    Object json.RawMessage
}



/*
//YS: when new transaction received from user, propagate it to other nodes.
func propagate_tx(the_tx Transaction) {

	//YS: this should be updated to use dynamic ip pool first
	ip_pool := strings.Split( IP_POOL , "_")

	for i := range ip_pool{
		//YS: do not dial itself
		if (ip_pool[i] == lan_ip) {
			continue
		}
		conn, err := net.Dial("tcp", ip_pool[i] + ":" + TCP_PORT)
		if err != nil {
			fmt.Println("CONNECTION ERROR:", err.Error())
			continue
		}

	//	json_data := json_transaction{ Data_type: "T", Tx: the_tx }

	//	encoder := json.NewEncoder(conn)
	//	encoder.Encode(json_data)
	}
}
*/
/*
//YS: Old simulation code

func handleConn(conn net.Conn) {
	defer conn.Close()
	//YS: current node write a message to the node that connect to it
	io.WriteString(conn, "Enter a new transaction:")
	scanner := bufio.NewScanner(conn)
	// take in transactions from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			// YS: get user input from the other node, and create []Transaction
			input_data := string(scanner.Text())
			t := time.Now()
			transaction_data := []Transaction {	Transaction{ TransactionId: input_data, Timestamp: t.String()}}

			//YS: using the transaction received from other node to generate block
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], transaction_data)
			if err != nil {
				panic (err)
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}
			//YS: This is to throw the new blockchain into the channel.
			//YS: should it be a json obj?? What is happening here?
			bcServer <- Blockchain
			io.WriteString(conn, "\nEnter a new transaction:")
		}
	}()
	// simulate receiving broadcast
	go func() {
		for {
			time.Sleep(30 * time.Second)
			mutex.Lock()
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			mutex.Unlock()
			io.WriteString(conn, string(output))
		}
	}()
	for _ = range bcServer {
		spew.Dump(Blockchain)
	}
}

*/
/*
func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}
*/
