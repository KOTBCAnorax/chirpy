package main

import (
	"encoding/json"
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
	w.Header().Add("content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	responseText := fmt.Sprintf(
		`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(responseText))
}

func (cfg *apiConfig) resetNumberOfHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	oldNumberOfHits := cfg.fileserverHits.Swap(0)
	responseText := fmt.Sprintf("Number of hits reset from %d to %d", oldNumberOfHits, cfg.fileserverHits.Load())
	w.Write([]byte(responseText))
}

type chirp struct {
	Body string `json:"body"`
}

type validChirpResponse struct {
	Valid bool `json:"valid"`
}

type errorChirpResponse struct {
	ErrorMsg string `json:"error"`
}

func generateErrorResponse(msg ...string) []byte {
	var errormsg string
	if len(msg) > 0 {
		errormsg = msg[0]
	} else {
		errormsg = "Something went wrong"
	}

	responseBody := errorChirpResponse{ErrorMsg: errormsg}
	data, _ := json.Marshal(responseBody)
	return data
}

func chirpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		w.Write(generateErrorResponse())
		return
	}

	if len(params.Body) > 140 {
		msg := "Chirp is too long"
		log.Print(msg)
		w.WriteHeader(400)
		w.Write(generateErrorResponse(msg))
		return
	}

	responseBody := validChirpResponse{Valid: true}
	data, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		w.WriteHeader(500)
		w.Write(generateErrorResponse())
		return
	}

	w.WriteHeader(200)
	w.Write(data)
}

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
