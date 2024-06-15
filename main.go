package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
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

func NewItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not create: %v", err)
		w.Write([]byte("could not create item"))
		return
	}
	h := md5.New()
	io.WriteString(h, item.Price)
	item.ID = fmt.Sprintf("%x", h.Sum(nil))

	res, err := json.MarshalIndent(item, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal: %v", err)
		w.Write([]byte("could not create item(2)"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getBlocks).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newItem).Methods("POST")

	log.Println("Running on port: 3000")
	log.Fatal(http.ListenAndServe(":3000", r))

}
