package main

import (
	"log"
	"net/http"

	"github.com/yimialmonte/GoAPI/handlers"
)

func main() {
	http.HandleFunc("/users/", handlers.UsersRouter)
	http.HandleFunc("/users", handlers.UsersRouter)
	http.HandleFunc("/", handlers.RootHandler)
	err := http.ListenAndServe("localhost:11111", nil)
	if err != nil {
		log.Fatal(err)
	}
}
