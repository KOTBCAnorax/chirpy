package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 Forbidden"))
		return
	}

	err := cfg.db.Reset(r.Context())
	if err != nil {
		log.Printf("Failed to delete users data: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(generateErrorResponse())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users data has been deleted\n"))
}
