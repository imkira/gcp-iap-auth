package jwt

import (
	"errors"
	"regexp"
)

// Config specifies the parameters for which to perform validation of JWT
// tokens in requests against.
type Config struct {
	PublicKeys     map[string]PublicKey
	MatchAudiences *regexp.Regexp
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
