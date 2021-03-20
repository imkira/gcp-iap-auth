package main

import (
	"log"
	"net"
	"net/http"
)

var (
	version  = "master"
	revision = "dev"
)

func main() {
	log.SetFlags(log.Flags() | log.LUTC)
	if len(revision) > 8 {
		revision = revision[:8]
	}
	log.Printf("Cloud IAP Auth & Proxy Server (build: %s.%s)\n", version, revision)

	if err := initConfig(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Matching audiences: %s\n", cfg.MatchAudiences)

	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/healthz", healthzHandler)

	if backend != nil && *backend != "" {
		proxy, err := newProxy(*backend, *emailHeader, timeout)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Proxying authenticated requests to backend %s", *backend)
		http.HandleFunc("/", proxy.handler)
	}

	addr := net.JoinHostPort(*listenAddr, *listenPort)
	if len(*tlsCertPath) != 0 || len(*tlsKeyPath) != 0 {
		listenAndServeHTTPS(addr)
	} else {
		listenAndServeHTTP(addr)
	}
}

func listenAndServeHTTP(addr string) {
	log.Printf("Listening on http://%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to listen on http://%s (%v)\n", addr, err)
	}
}

func listenAndServeHTTPS(addr string) {
	log.Printf("Listening on https://%s\n", addr)
	err := http.ListenAndServeTLS(addr, *tlsCertPath, *tlsKeyPath, nil)
	if err != nil {
		log.Fatalf("Failed to listen on https://%s (%v)\n", addr, err)
	}
}
