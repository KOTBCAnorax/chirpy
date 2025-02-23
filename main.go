package main

import (
	"log"
	"net/http"
)

func main() {
	const filePathRoot = "."
	const port = "8080"
	var apiCfg = apiConfig{}
	mux := http.NewServeMux()
	fileServerHandler := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/", (http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServerHandler))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.numberOfHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetNumberOfHits)
	mux.HandleFunc("POST /api/validate_chirp", chirpHandler)

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
