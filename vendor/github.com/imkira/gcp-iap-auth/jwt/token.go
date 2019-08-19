package jwt

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

const (
	tokenHeader    = "X-Goog-IAP-JWT-Assertion"
	algorithm      = "ES256"
	algorithmClaim = "alg"
	keyIDClaim     = "kid"
	issuerClaim    = "https://cloud.google.com/iap"
)

func tokenKey(token *jwt.Token) (interface{}, error) {
	if _, ok := tokenMethod(token); !ok {
		return nil, fmt.Errorf("Invalid algorithm: %v", token.Header[algorithmClaim])
	}
	keyID, _ := token.Header[keyIDClaim].(string)
	key := token.Claims.(*Claims).cfg.PublicKeys[keyID]
	if len(key) == 0 {
		return nil, fmt.Errorf("No public key for %q", keyID)
	}
	parsedKey, err := jwt.ParseECPublicKeyFromPEM(key)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse key: %v", err)
	}
	return parsedKey, nil
}

func tokenMethod(token *jwt.Token) (jwt.SigningMethod, bool) {
	if token.Header[algorithmClaim] != algorithm {
		return nil, false
	}
	method, ok := token.Method.(*jwt.SigningMethodECDSA)
	if !ok {
		return nil, false
	}
	return method, true
}
