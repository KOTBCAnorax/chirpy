package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type chirp struct {
	Body string `json:"body"`
}

type validChirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func generateErrorResponse(w http.ResponseWriter, statusCode int, msg ...string) {
	var errormsg string
	if len(msg) > 0 {
		errormsg = msg[0]
	} else {
		errormsg = "Something went wrong\n"
	}

	w.Header().Set("content_type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(errormsg))
}

func generateValidResponse(w http.ResponseWriter, chirp string) []byte {
	filteredChirp := FilterProfane(chirp)
	responseBody := validChirpResponse{CleanedBody: filteredChirp}
	data, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error encoding response: %s\n", err)
		generateErrorResponse(w, 500)
		return nil
	}

	return data
}

func chirpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		generateErrorResponse(w, 500)
		return
	}

	if len(params.Body) > 140 {
		msg := "Chirp is too long"
		log.Println(msg)
		generateErrorResponse(w, 400, msg)
		return
	}

	data := generateValidResponse(w, params.Body)

	w.WriteHeader(200)
	w.Write(data)
}
