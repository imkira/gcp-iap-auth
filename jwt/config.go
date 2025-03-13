package jwt

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	jwt "github.com/dgrijalva/jwt-go"
)

// Config specifies the parameters for which to perform validation of JWT
// tokens in requests against.
type Config struct {
	PublicKeyPath  string
	PublicKeys     map[string]PublicKey
	MatchAudiences *regexp.Regexp
	MatchDomains   map[string]bool
}

func (cfg *Config) RefreshPublicKeys() error {
	var err error
	if len(cfg.PublicKeyPath) != 0 {
		cfg.PublicKeys, err = loadPublicKeysFromFile(cfg.PublicKeyPath)
	} else {
		cfg.PublicKeys, err = FetchPublicKeys()
	}
	return err
}

func (cfg *Config) GetPublicKey(keyID string) (interface{}, error) {
	key, ok := cfg.PublicKeys[keyID]
	// Refresh public keys on lookup failure. For details, see
	// https://github.com/imkira/gcp-iap-auth/issues/10 and
	// https://stackoverflow.com/questions/44828856/google-iap-public-keys-expiry.
	if !ok {
		if err := cfg.RefreshPublicKeys(); err != nil {
			return nil, err
		}
		key, ok = cfg.PublicKeys[keyID]
		if !ok {
			return nil, fmt.Errorf("No public key for %q", keyID)
		}
	}
	parsedKey, err := jwt.ParseECPublicKeyFromPEM(key)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse key: %v", err)
	}
	return parsedKey, nil
}

// Validate validates the Configuration.
func (cfg *Config) Validate() error {
	if cfg.MatchAudiences == nil {
		return errors.New("No audiences to match defined")
	}
	if len(cfg.PublicKeys) == 0 {
		return errors.New("No public keys defined")
	}
	return nil
}

func (cfg *Config) matchesAudience(aud *Audience) bool {
	return cfg.MatchAudiences.MatchString((string)(*aud))
}

func (cfg *Config) matchesDomain(hd string) bool {
	return len(cfg.MatchDomains) == 0 || cfg.MatchDomains[hd]
}

func loadPublicKeysFromFile(filePath string) (map[string]PublicKey, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return DecodePublicKeys(f)
}
