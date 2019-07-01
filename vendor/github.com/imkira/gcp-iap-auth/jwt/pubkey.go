package jwt

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

var (
	// PublicKeysURL is the URL from which public keys will be fetched.
	PublicKeysURL = "https://www.gstatic.com/iap/verify/public_key"
	// HTTPClient is the default HTTP Client to use for fetching public keys.
	HTTPClient = &http.Client{Timeout: 10 * time.Second}
)

// PublicKey are Google's public keys to use for JWT token validation.
type PublicKey []byte

// CreatePublicKey creates a PublicKey from a byte slice.
func CreatePublicKey(b []byte) PublicKey {
	return PublicKey(b)
}

// DecodePublicKeys decodes all public keys from the given Reader.
func DecodePublicKeys(r io.Reader) (map[string]PublicKey, error) {
	var skeys map[string]string
	if err := json.NewDecoder(r).Decode(&skeys); err != nil {
		return nil, err
	}
	bkeys := make(map[string]PublicKey)
	for k, v := range skeys {
		if len(v) != 0 {
			bkeys[k] = CreatePublicKey([]byte(v))
		}
	}
	return bkeys, nil
}

// FetchPublicKeys downloads and decodes all public keys from Google.
func FetchPublicKeys() (map[string]PublicKey, error) {
	r, err := HTTPClient.Get(PublicKeysURL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return DecodePublicKeys(r.Body)
}
