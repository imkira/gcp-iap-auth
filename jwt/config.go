package jwt

import "errors"

// Config specifies the parameters for which to perform validation of JWT
// tokens in requests against.
type Config struct {
	PublicKeys map[string]PublicKey
	Audiences  []*Audience
}

// Validate validates the Configuration.
func (cfg *Config) Validate() error {
	if len(cfg.Audiences) == 0 {
		return errors.New("No audiences defined")
	}
	for _, aud := range cfg.Audiences {
		if err := aud.Validate(); err != nil {
			return err
		}
	}
	if len(cfg.PublicKeys) == 0 {
		return errors.New("No public keys defined")
	}
	return nil
}

func (cfg *Config) containsAudience(aud *Audience) bool {
	for _, aud2 := range cfg.Audiences {
		if *aud == *aud2 {
			return true
		}
	}
	return false
}
