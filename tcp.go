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
)


// bcServer handles incoming concurrent Blocks
//var bcServer chan []Block


var	ext_ip = GetOutboundIP()
var	lan_ip = GetLAN_IP()

func Runtcp() error {

	//bcServer = make(chan []Block)
	
	server, err := net.Listen("tcp", ":" + os.Getenv("ADDRTCP"))
	
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	
	//YS: Print Node External and Internal IP	

	fmt.Println("My External IP: " + ext_ip );
	fmt.Println("My Internal IP: " + lan_ip );		
	
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

	ip_pool := strings.Split( os.Getenv("IP_POOL") , "_")
	broadcast_invl, err :=  strconv.Atoi(os.Getenv("BC_BROADCAST_INVL"))
	if err != nil {
        fmt.Println("broadcast_invl error, default value = 600 loaded. ", err.Error())
        broadcast_invl = 600
    }
	
	for v := range ip_pool{
		fmt.Println("ip_pool : " + ip_pool[v])
	}
	
	for {		
		for i := range ip_pool{
			conn, err := net.Dial("tcp", ip_pool[i] + ":" + os.Getenv("ADDRTCP"))
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
YS 4/20:  
Accept TCP connection
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
	
	//YS: in blockchain.go, replace with new chain if it is longer
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
			
			transaction_data := []Transaction {	Transaction{ transactionId: input_data, Timestamp: t.String()}}
			
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


