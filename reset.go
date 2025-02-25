package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		log.Println("Unauthorized attempt to delete all users data")
		generateErrorResponse(w, 403, "403 Forbidden")
		return
	}

	err := cfg.db.Reset(r.Context())
	if err != nil {
		log.Printf("Failed to delete users data: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users data has been deleted\n"))
}
