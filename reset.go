package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits_msg := fmt.Sprintf("File Server Hits Reset to: %v", cfg.fileserverHits.Load())
	w.Write([]byte(hits_msg))
}
