package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/imkira/gcp-iap-auth/jwt"
	"github.com/namsral/flag"
)

const flagEnvPrefix = "GCP_IAP_AUTH"

var (
	cfg            = &jwt.Config{}
	listenAddr     = flag.String("listen-addr", "0.0.0.0", "Listen address")
	listenPort     = flag.String("listen-port", "", "Listen port (default: 80 for HTTP or 443 for HTTPS)")
	audiences      = flag.String("audiences", "", "Comma-separated list of JWT Audiences (elements can be paths like \"/projects/PROJECT_NUMBER/apps/PROJECT_ID\" or regular expressions like \"/^\\/projects\\/PROJECT_NUMBER/.*\" if you enclose them in slashes)")
	publicKeysPath = flag.String("public-keys", "", "Path to public keys file (optional)")
	tlsCertPath    = flag.String("tls-cert", "", "Path to TLS server's, intermediate's and CA's PEM certificate (optional)")
	tlsKeyPath     = flag.String("tls-key", "", "Path to TLS server's PEM key file (optional)")
	backend        = flag.String("backend", "", "Proxy authenticated requests to the specified URL (optional)")
	emailHeader    = flag.String("email-header", "X-WEBAUTH-USER", "In proxy mode, set the authenticated email address in the specified header")
	duration       = flag.String("timeout", "30s", "proxy request timeout, e.g. 15s or 1m")
	timeout        time.Duration
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
	if err := initAudiences(*audiences); err != nil {
		return err
	}
	if err := initPublicKeys(*publicKeysPath); err != nil {
		return err
	}
	if err := initTimeout(*duration); err != nil {
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

func initAudiences(audiences string) error {
	str, err := extractAudiencesRegexp(audiences)
	if err != nil {
		return err
	}
	re, err := regexp.Compile(str)
	if err != nil {
		return fmt.Errorf("Invalid audiences regular expression %q (%v)", str, err)
	}
	cfg.MatchAudiences = re
	return nil
}

func extractAudiencesRegexp(audiences string) (string, error) {
	var strs []string
	for _, audience := range strings.Split(audiences, ",") {
		str, err := extractAudienceRegexp(audience)
		if err != nil {
			return "", err
		}
		strs = append(strs, str)
	}
	return strings.Join(strs, "|"), nil
}

func extractAudienceRegexp(audience string) (string, error) {
	if strings.HasPrefix(audience, "/") && strings.HasSuffix(audience, "/") {
		if len(audience) < 3 {
			return "", fmt.Errorf("Invalid audiences regular expression %q", audience)
		}
		return audience[1 : len(audience)-1], nil
	}
	return parseRawAudience(audience)
}

func parseRawAudience(audience string) (string, error) {
	aud, err := jwt.ParseAudience(audience)
	if err != nil {
		return "", fmt.Errorf("Invalid audience %q (%v)", audience, err)
	}
	return fmt.Sprintf("^%s$", regexp.QuoteMeta((string)(*aud))), nil
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

func initTimeout(duration string) error {
	var err error
	timeout, err = time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("failed to parse timeout (%s): %v", duration, err)
	}
	return nil
}

func loadPublicKeysFromFile(filePath string) (map[string]jwt.PublicKey, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return jwt.DecodePublicKeys(f)
}
