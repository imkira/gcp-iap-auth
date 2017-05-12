package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/imkira/gcp-iap-auth/jwt"
)

type proxy struct {
	backend     *url.URL
	emailHeader string
	proxy       *httputil.ReverseProxy
}

func newProxy(backendURL, emailHeader string) (*proxy, error) {
	backend, err := url.Parse(backendURL)
	if err != nil {
		return nil, fmt.Errorf("Could not parse URL '%s': %s", backendURL, err)
	}
	return &proxy{
		backend:     backend,
		emailHeader: emailHeader,
		proxy:       httputil.NewSingleHostReverseProxy(backend),
	}, nil
}

func (p *proxy) handler(res http.ResponseWriter, req *http.Request) {
	claims, err := jwt.RequestClaims(req, cfg)
	if err != nil {
		if claims == nil || len(claims.Email) == 0 {
			log.Printf("Failed to authenticate (%v)\n", err)
		} else {
			log.Printf("Failed to authenticate %q (%v)\n", claims.Email, err)
		}
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if p.emailHeader != "" {
		req.Header.Set(p.emailHeader, claims.Email)
	}
	p.proxy.ServeHTTP(res, req)
}
