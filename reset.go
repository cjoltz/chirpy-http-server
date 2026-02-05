package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// Clear DB
	cfg.db.DeleteUsers(r.Context())
	// Reset Server Hits Count
	cfg.fileserverHits.Store(0)
	// Write Header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server Reset"))
}
