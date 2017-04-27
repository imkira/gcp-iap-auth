package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/imkira/gcp-iap-auth/jwt"
	"github.com/namsral/flag"
)

const flagEnvPrefix = "GCP_IAP_AUTH"

var (
	cfg            = &jwt.Config{}
	listenAddr     = flag.String("listen-addr", "0.0.0.0", "Listen address")
	listenPort     = flag.String("listen-port", "", "Listen port (default: 80 for HTTP or 443 for HTTPS)")
	audiences      = flag.String("audiences", "", "Comma separated list of JWT Audiences (format: https://yourdomain or https://yourdomain:port)")
	publicKeysPath = flag.String("public-keys", "", "Path to public keys file (optional)")
	tlsCertPath    = flag.String("tls-cert", "", "Path to TLS server's, intermediate's and CA's PEM certificate (optional)")
	tlsKeyPath     = flag.String("tls-key", "", "Path to TLS server's PEM key file (optional)")
)

func initConfig() error {
	flag.EnvironmentPrefix = flagEnvPrefix
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)
	flag.Parse()
	if err := initServerPort(); err != nil {
		return err
	}
	if len(*audiences) == 0 {
		return errors.New("You must specify --audiences")
	}
	if err := initAudiences(strings.Split(*audiences, ",")); err != nil {
		return err
	}
	if err := initPublicKeys(*publicKeysPath); err != nil {
		return err
	}
	return nil
}

func initServerPort() error {
	if len(*listenPort) == 0 {
		if len(*tlsCertPath) != 0 || len(*tlsKeyPath) != 0 {
			*listenPort = "443"
		} else {
			*listenPort = "80"
		}
	}
	if _, err := strconv.Atoi(*listenPort); err != nil {
		return fmt.Errorf("Invalid listen port %q", *listenPort)
	}
	return nil
}

func initAudiences(rawURLs []string) error {
	for _, rawURL := range rawURLs {
		aud, err := jwt.ParseAudience(rawURL)
		if err != nil {
			return fmt.Errorf("Invalid audience %q (%v)", rawURL, err)
		}
		cfg.Audiences = append(cfg.Audiences, aud)
	}
	return nil
}

func initPublicKeys(filePath string) error {
	var err error
	if len(filePath) != 0 {
		cfg.PublicKeys, err = loadPublicKeysFromFile(filePath)
	} else {
		cfg.PublicKeys, err = jwt.FetchPublicKeys()
	}
	if err != nil {
		return err
	}
	return cfg.Validate()
}

func loadPublicKeysFromFile(filePath string) (map[string]jwt.PublicKey, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return jwt.DecodePublicKeys(f)
}
