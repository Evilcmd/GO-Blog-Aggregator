package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJson(res http.ResponseWriter, code int, payload interface{}) {
	res.WriteHeader(code)
	dat, _ := json.Marshal(payload)
	res.Write(dat)
}

func respondWithError(res http.ResponseWriter, code int, message string) {
	payload := struct {
		Error string `json:"error"`
	}{
		message,
	}
	respondWithJson(res, code, &payload)
}

func healthz(res http.ResponseWriter, req *http.Request) {
	payload := struct {
		Status string `json:"status"`
	}{
		"ok",
	}
	respondWithJson(res, 200, &payload)
}

func erreHandle(res http.ResponseWriter, req *http.Request) {
	respondWithError(res, 200, "Internal Server Error")
}
