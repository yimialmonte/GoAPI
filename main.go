package main

import (
	"log"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Running API v1\n"))
}

func main() {
	http.HandleFunc("/", rootHandler)
	err := http.ListenAndServe("localhost:11111", nil)
	if err != nil {
		log.Fatal(err)
	}
}
