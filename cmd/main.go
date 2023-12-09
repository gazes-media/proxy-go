package main

import (
	"log"
	"net/http"
	"os"

	"github.com/trail-l31/gazes-proxy/internal"
)

var port string

func init() {
	if port = os.Getenv("PORT"); port == "" {
		log.Fatal("Missing env variable 'PORT'")
	}
}

func main() {
	http.HandleFunc("/", internal.ProxyHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
