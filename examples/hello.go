package main

import (
	"io"
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
// our server at http://127.0.0.1:12345/hello
// For valid tokens we say Hello, for others we say Sorry.
func HelloHandler(w http.ResponseWriter, req *http.Request) {
	claims, err := jwt.RequestClaims(req, cfg)
	if err != nil {
		if claims == nil || len(claims.Email) == 0 {
			io.WriteString(w, "Sorry: "+err.Error())
		} else {
			io.WriteString(w, "Sorry "+claims.Email+":"+err.Error())
		}
	} else {
		io.WriteString(w, "Hello "+claims.Email)
	}
}

func main() {
	http.HandleFunc("/hello", HelloHandler)
	addr := "127.0.0.1:12345"
	log.Printf("Running at http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
