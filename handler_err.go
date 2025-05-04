package main

import (
	"net/http"
)

func errHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something went wrong")
}
