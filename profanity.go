package main

import "strings"

var profaneWords = map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}

func FilterProfane(chirp string) string {
	words := strings.Split(chirp, " ")
	lowerChirp := strings.ToLower(chirp)
	lowerWords := strings.Split(lowerChirp, " ")
	for i, word := range lowerWords {
		if profaneWords[word] {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
