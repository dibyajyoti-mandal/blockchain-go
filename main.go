package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Block struct {
	Index    int
	Data     Checkout
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

type Checkout struct {
	ItemID    string `json:"item_id"`
	Buyer     string `json:"buyer"`
	Date      string `json:"date"`
	IsGenesis bool   `json:"is_genesis"`
}

type Blockchain struct {
	blocks []*Block
}

var BlockChain *Blockchain

func CreateBlock(prev *Block, checkout Checkout) *Block {

	block := &Block{}
	block.Index = prev.Index + 1
	block.Time = time.Now().String()

	return block
}

func (bc *Blockchain) addBlock(data Checkout) {
	prev := bc.blocks[len(bc.blocks)-1]
	block := CreateBlock(prev, data)

	if valid(block, prev) {
		bc.blocks = append(bc.blocks, block)
	}
}

func newItem(w http.ResponseWriter, r *http.Request) {
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

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkout Checkout
	if err := json.NewDecoder(r.Body).Decode(&checkout); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("coudlnt create block")
		w.Write([]byte("COuld not create block"))
		return
	}

	BlockChain.addBlock(checkout)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getBlocks).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newItem).Methods("POST")

	log.Println("Running on port: 3000")
	log.Fatal(http.ListenAndServe(":3000", r))

}
