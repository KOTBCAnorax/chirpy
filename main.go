package main

import (
	"log"
	"net/http"
)

func main() {
	const filePathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filePathRoot)))

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Fatal(srv.ListenAndServe())
}
