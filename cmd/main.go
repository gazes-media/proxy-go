package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"proxy/internal"
)

var port string

func init() {
	if port = os.Getenv("PORT"); port == "" {
		log.Fatal("Missing env variable 'PORT'")
	}
}

func main() {
	http.HandleFunc("/", internal.HandleIndex)
	serve(port)
}

// The serve function starts a web server on the specified port and logs any errors that occur.
func serve(port string) {
	fmt.Printf("Server started on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
