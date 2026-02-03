package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello from my Chirpy back-end...")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}