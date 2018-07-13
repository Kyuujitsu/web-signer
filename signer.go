package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewSignServer(&btcSigner{})

	banner := `

Usage:
curl "localhost:3434/sign?data=[datahex]&private=[privatehex]"

Both data and private key are expected to be hex-encoded.

Server listening localhost port 3434...
`
	log.Print(banner)

	if err := http.ListenAndServe(":3434", server); err != nil {
		log.Fatalf("could not listen on port 3434: %v", err)
	}
}
