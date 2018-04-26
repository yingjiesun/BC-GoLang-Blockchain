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
	//"io"
	"log"
	"net"
	"os"
	"time"
	"fmt"
	"strings"
	"github.com/davecgh/go-spew/spew"
	"strconv"
	"bytes"
  "io/ioutil"
  "net/http"
	"math/rand"
)

// bcServer handles incoming concurrent Blocks
//var bcServer chan []Block

var	ext_ip = GetOutboundIP()
var	lan_ip = GetLAN_IP()

func Runtcp() error {

	//bcServer = make(chan []Block)

	server, err := net.Listen("tcp", ":" + TCP_PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	//YS: Print Node External and Internal IP
	fmt.Println("My External IP: " + ext_ip );
	fmt.Println("My Internal IP: " + lan_ip );

	/*
	YS: Testing code below, this is for testing to generate block randomly in 30-90 seconds
	*/
	go func(){
		for {
			fmt.Println("Start creating new block")
			if len(temp_trans) > 0 {
				mutex.Lock()
				newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], temp_trans)				
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
					spew.Dump(Blockchain)
				} else {
					fmt.Println("INVALID block")
				}
			}
			time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
		}
	}()

	//YS: END of generating test block

	/*
	YS: Testing code below, this is for testing to generate transaction randomly in 10-20 seconds
	*/
	go func(){
		test_tran_id := 100
		for {
			t := time.Now()
			var tranaction_new = Transaction{ TransactionId: strconv.Itoa(test_tran_id), Timestamp: t.String()}
			mutex.Lock()
			append_temp_trans(tranaction_new)
			mutex.Unlock()
			test_tran_id++
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		}
	}()

	//YS: END of generating test transaction

	go dialConn()

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
	fmt.Printf("outer IP: " + localAddr[0:idx] );
    return localAddr[0:idx]
}

/*
YS 4/20:
Create TCP connection to other nodes listed in IP_POOL
Every 30 seconds, dial TCP and send blockchain
*/

func dialConn() {
	ip_pool := strings.Split( IP_POOL , "_")
	broadcast_invl := BROADCAST_INTERVAL
	for v := range ip_pool{
		fmt.Println("ip_pool : " + ip_pool[v])
	}
	for {
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
			encoder := json.NewEncoder(conn)
			encoder.Encode(Blockchain)
		}
		time.Sleep(time.Duration(broadcast_invl) * time.Second)
	}
}


/*
YS: Accept TCP connection
*/

func handleConn(conn net.Conn) {
	var receivedBlockchain []Block
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	decoder.Decode(&receivedBlockchain)
	t := time.Now()
	fmt.Println("===============================================")
	fmt.Println("My External IP: " + ext_ip );
	fmt.Println("My Internal IP: " + lan_ip );
	fmt.Println("Received blockchain: " + t.String())
	spew.Dump(receivedBlockchain)
	//YS: function in blockchain.go, replace with new chain if it is longer
	replaceChain(receivedBlockchain)
  conn.Close()
}

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

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}
