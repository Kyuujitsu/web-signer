package main

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/btcsuite/btcd/btcec"
)

// Signer will sign a data with private key provided
type Signer interface {
	Sign(datahex, privhex string) (string, error)
}

type btcSigner struct{}

// Sign returns signature transaction hash for provided data and private key
func (s *btcSigner) Sign(datahex, privhex string) (string, error) {
	priv, err := hex.DecodeString(privhex)
	if err != nil {
		return "", err
	}

	data, err := hex.DecodeString(datahex)
	if err != nil {
		return "", err
	}

	key, _ := btcec.PrivKeyFromBytes(btcec.S256(), priv)
	sig, err := key.Sign(data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig.Serialize()), nil
}

// SignServer responds to GET request on /sign and returns signature for
// provided data
type SignServer struct {
	signer Signer
	http.Handler
}

// NewSignServer creates SignServer with routing
func NewSignServer(signer Signer) *SignServer {
	s := new(SignServer)
	s.signer = signer

	router := http.NewServeMux()
	router.Handle("/sign", http.HandlerFunc(s.signHandler))

	s.Handler = router

	return s
}

func (s *SignServer) signHandler(w http.ResponseWriter, r *http.Request) {
	data, ok := r.URL.Query()["data"]
	if !ok || len(data[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	keys, ok := r.URL.Query()["private"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sig, err := s.signer.Sign(data[0], keys[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, sig)
}
