package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/wedaly/appmeta/server"
)

var addressFlag = flag.String("a", ":8000", "Listen address for the server")

func main() {
	flag.Parse()
	log.Printf("Starting server on %s\n", *addressFlag)
	err := http.ListenAndServe(*addressFlag, server.Handler())
	if err != nil {
		log.Fatalf("Unexpected error occurred: %s\n", err)
	}
}
