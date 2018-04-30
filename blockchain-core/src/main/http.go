package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	//"os"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"net/url"
	"time"
)

func run() error {
	mux := makeMuxRouter()
	httpAddr := HTTP_PORT
	log.Println("Listening on ", HTTP_PORT)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET", "OPTIONS")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	//YS: Enable CORS
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var bytes []byte
	var err error

	r.ParseForm()
	preBlockHash := retrieveParameterFromRequest(r.Form, "offsetByPreBlockHash")

	if preBlockHash != "" {
		searchedBlockIndex := searchSpecifiedBlockByHash(preBlockHash)
		if searchedBlockIndex == -1 {
			bytes = []byte{}
		} else {
			bytes, err = json.MarshalIndent(Blockchain[searchedBlockIndex+1:], "", "  ")
		}
	} else {
		bytes, err = json.MarshalIndent(Blockchain, "", "  ")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func searchSpecifiedBlockByHash(hashCode string) int {
	for index, block := range Blockchain {
		if block.Hash == hashCode {
			return index
		}
	}
	return -1
}

func retrieveParameterFromRequest(values url.Values, key string) string {
	v := values.Get(key)
	if v == "" {
		return ""
	} else {
		return v
	}
}

type Message struct {
	transactions []Transaction
}

/*
YS:
TODO: modify this function to take jason obj and convert to []Transaction
For now a new block generated at same time. In future new block is triggered by time.
*/

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], m.transactions, production_ip)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
		spew.Dump(Blockchain)
	}
	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
