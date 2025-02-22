package main

import (
	"fmt"
	"net/http"
)

func main() {
	serverHandler := http.NewServeMux()
	server := http.Server{Handler: serverHandler, Addr: ":8080"}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
