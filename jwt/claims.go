package jwt

import (
	"fmt"

	jwt "github.com/golang-jwt/jwt/v4"
)

// Claims represents parsed JWT Token Claims.
type Claims struct {
	jwt.StandardClaims
	Domain string `json:"hd,omitempty"`
	Email  string `json:"email,omitempty"`
	cfg    *Config
}

// Valid validates the Claims.
func (c Claims) Valid() error {
	if err := (c.StandardClaims).Valid(); err != nil {
		return err
	}
	if c.Issuer != issuerClaim {
		return fmt.Errorf("Invalid issuer: %q", c.Issuer)
	}
	aud, err := ParseAudience(c.Audience)
	if err != nil {
		return fmt.Errorf("Invalid audience %q: %v", c.Audience, err)
	}
	if !c.cfg.matchesAudience(aud) {
		return fmt.Errorf("Unexpected audience: %q", c.Audience)
	}
	if !c.cfg.matchesDomain(c.Domain) {
		return fmt.Errorf("Unexpected domain: %q", c.Domain)
	}
	return nil
}
