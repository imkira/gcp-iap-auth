package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
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
	cfg := token.Claims.(*Claims).cfg
	return cfg.GetPublicKey(keyID)
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
