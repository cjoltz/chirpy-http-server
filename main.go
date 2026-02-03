package main

import (
	"log"
	"net/http"
)

func main() {
	const filePathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	// App
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	mux.HandleFunc("/healthz", handlerReadiness)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}


func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}