package main

import (
	"log"

	"net/http"

	"github.com/gorilla/mux"
)

type Block struct {
	Pos      int
	Data     ItemCheckout
	Time     string
	Hash     string
	PrevHash string
}

type Item struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Seller string `json:"seller"`
	Price  string `json:"price"`
}

type ItemCheckout struct {
	ItemID    string `json:"item_id"`
	Buyer     string `json:"buyer"`
	Date      string `json:"date"`
	IsGenesis bool   `json:"is_genesis"`
}

type Blockchain struct {
	blocks []*Block
}

var BlockChain *Blockchain

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getBlocks).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newItem).Methods("POST")

	log.Println("Running on port: 3000")
	log.Fatal(http.ListenAndServe(":3000", r))

}
