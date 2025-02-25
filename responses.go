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

type errorChirpResponse struct {
	ErrorMsg string `json:"error"`
}

func generateErrorResponse(msg ...string) []byte {
	var errormsg string
	if len(msg) > 0 {
		errormsg = msg[0]
	} else {
		errormsg = "Something went wrong\n"
	}

	responseBody := errorChirpResponse{ErrorMsg: errormsg}
	data, _ := json.Marshal(responseBody)
	return data
}

func generateValidResponse(w http.ResponseWriter, chirp string) []byte {
	filteredChirp := FilterProfane(chirp)
	responseBody := validChirpResponse{CleanedBody: filteredChirp}
	data, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error encoding response: %s\n", err)
		w.WriteHeader(500)
		w.Write(generateErrorResponse())
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

	data := generateValidResponse(w, params.Body)

	w.WriteHeader(200)
	w.Write(data)
}
