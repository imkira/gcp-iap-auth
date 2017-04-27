package main

import (
	"log"
	"net/http"

	"github.com/imkira/gcp-iap-auth/jwt"
)

var cfg *jwt.Config

// In here we initialize the configuration for our app.
// It doesn't need to be in "init".
func init() {
	audience, err := jwt.ParseAudience("http://your.domain.com")
	if err != nil {
		log.Fatal(err)
	}
	publicKeys, err := jwt.FetchPublicKeys()
	if err != nil {
		log.Fatal(err)
	}
	cfg = &jwt.Config{
		Audiences:  []*jwt.Audience{audience},
		PublicKeys: publicKeys,
	}
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
}

// Here we validate the tokens in all requests going to
// our server at http://127.0.0.1:12345/auth
// For valid tokens we return 200, otherwise 401.
func AuthHandler(w http.ResponseWriter, req *http.Request) {
	if err := jwt.ValidateRequestClaims(req, cfg); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	http.HandleFunc("/auth", AuthHandler)
	addr := "127.0.0.1:12345"
	log.Printf("Running at http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
