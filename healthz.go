package main

import (
	"fmt"
	"net/http"
)

func healthzHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	fmt.Fprintln(res, `{"status":"ok"}`)
}
