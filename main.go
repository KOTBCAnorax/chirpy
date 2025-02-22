package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) numberOfHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	responseText := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(responseText))
}

func (cfg *apiConfig) resetNumberOfHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	oldNumberOfHits := cfg.fileserverHits.Swap(0)
	responseText := fmt.Sprintf("Number of hits reset from %d to %d", oldNumberOfHits, cfg.fileserverHits.Load())
	w.Write([]byte(responseText))
}

func main() {
	const filePathRoot = "."
	const port = "8080"
	var apiCfg = apiConfig{}
	mux := http.NewServeMux()
	fileServerHandler := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/", (http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServerHandler))))

	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apiCfg.numberOfHits)
	mux.HandleFunc("POST /reset", apiCfg.resetNumberOfHits)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
