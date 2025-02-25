package main

import (
	"net/http"
)

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
