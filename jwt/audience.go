package jwt

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
)

// Audience must be the base URL from the request including protocol, domain,
// and port if applicable for the domains you specify in your IAP proxy. For
// example, https://example.com or https://foo.example.com:port.
type Audience url.URL

// NewAudience returns an Audience from a URL.
func NewAudience(u *url.URL) *Audience {
	return (*Audience)(u)
}

// Sanitize normalizes the structure of the Audience's URL and validates it.
func (aud *Audience) Sanitize() error {
	aud.Scheme = strings.ToLower(aud.Scheme)
	var err error
	var port string
	host := aud.Host
	if strings.LastIndex(host, ":") >= 0 {
		host, port, err = net.SplitHostPort(host)
		if err != nil {
			return err
		}
	}
	if len(port) == 0 {
		defaultPort := "80"
		if aud.Scheme == "https" {
			defaultPort = "443"
		}
		aud.Host = net.JoinHostPort(host, defaultPort)
	}
	return aud.Validate()
}

// Validate performs error checking on the Audience's URL.
func (aud *Audience) Validate() error {
	if aud.Scheme != "http" && aud.Scheme != "https" {
		return fmt.Errorf("Unexpected scheme: %s", aud.Scheme)
	}
	host, _, err := net.SplitHostPort(aud.Host)
	if err != nil {
		return err
	}
	if len(host) == 0 {
		return errors.New("Host not specified")
	}
	if aud.User != nil {
		return errors.New("Not expecting user")
	}
	if len(aud.Path) != 0 || len(aud.RawPath) != 0 {
		return errors.New("Not expecting path")
	}
	if len(aud.RawQuery) != 0 {
		return errors.New("Not expecting query")
	}
	if len(aud.Fragment) != 0 {
		return errors.New("Not expecting fragment")
	}
	return nil
}

// ParseAudience parses an Audience from a URL string.
func ParseAudience(rawURL string) (*Audience, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	aud := NewAudience(u)
	if err := aud.Sanitize(); err != nil {
		return nil, err
	}
	return aud, nil
}
