package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe("localhost:11111", nil)
	if err != nil {
		log.Fatal(err)
	}
}
