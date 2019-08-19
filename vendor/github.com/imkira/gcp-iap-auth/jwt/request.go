package jwt

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// ValidateRequestClaims checks the validity of the claims in the request.
func ValidateRequestClaims(req *http.Request, cfg *Config) error {
	_, err := RequestClaims(req, cfg)
	return err
}

// RequestClaims checks the validity and returns the claims in the request.
// Claims may be returned even if an error occurs.
func RequestClaims(req *http.Request, cfg *Config) (*Claims, error) {
	tokenString, err := tokenStringFromRequest(req)
	if err != nil {
		return nil, err
	}
	claims := &Claims{cfg: cfg}
	_, err = jwt.ParseWithClaims(tokenString, claims, tokenKey)
	return claims, err
}

func tokenStringFromRequest(req *http.Request) (string, error) {
	token := req.Header.Get(tokenHeader)
	if len(token) == 0 {
		return "", errors.New("Token was not found in the request headers")
	}
	return token, nil
}
