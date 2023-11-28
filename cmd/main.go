package main

import (
	"log"
	"net/http"

	"github.com/trail-l31/gazes-proxy/internal"
)

func main() {
	http.HandleFunc("/", internal.ProxyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
