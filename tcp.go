package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"	
	"sync"
	"time"
	"fmt"
	"strings"
	"github.com/davecgh/go-spew/spew"
	//"strconv"
)


// bcServer handles incoming concurrent Blocks
var bcServer chan []Block
var mutex = &sync.Mutex{}

func Runtcp() error {

	// start TCP and serve TCP server
	
	bcServer = make(chan []Block)
	
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDRTCP"))
	
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	
	GetOutboundIP()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		
		localAddr := conn.LocalAddr().String()
		idx := strings.LastIndex(localAddr, ":")
		fmt.Printf("outer IP: " + localAddr[0:idx] );			
		
		go handleConn(conn)
	}	
		
	return nil
}


//function to get the public ip address
//YS: currently this is only to display IP in termainl

func GetOutboundIP() string {
	
    conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
				log.Fatal(err)
		}
		
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	fmt.Printf("outer IP: " + localAddr[0:idx] );
    return localAddr[0:idx]

	/*
	YS: get real external IP
	
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
	*/
}

func handleConn(conn net.Conn) {

	defer conn.Close()

	io.WriteString(conn, "Enter a new transaction:")

	scanner := bufio.NewScanner(conn)

	// take in transactions from stdin and add it to blockchain after conducting necessary validation
	
	/*
	
	YS: old code using BPM
	
	go func() {
		for scanner.Scan() {
			
			bpm, err := strconv.Atoi(scanner.Text())
			
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], bpm)
			if err != nil {
				log.Println(err)
				continue
			}
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}

			bcServer <- Blockchain
			io.WriteString(conn, "\nEnter a new transaction:")
		}
	}()
	
	*/
	
	go func() {
		for scanner.Scan() {
						
			// YS: get user input and create []Transaction
			
			input_data := string(scanner.Text())
			
			var transaction_data = []Transaction {	Transaction{ transactionId: input_data	} }
			
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], transaction_data)
		
			if err != nil {
				panic (err)
			}
		
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}

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




