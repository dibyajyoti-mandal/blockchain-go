package main

import (
	"fmt"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("hello")
	r := mux.NewRouter()
	r.HandleFunc("/").Methods("GET")
	r.HandleFunc("/").Methods("POST")
	r.HandleFunc("/new", newItem).Methods("POST")

}
