package jwt

import (
)

// From the IAP docs at https://cloud.google.com/iap/docs/signed-headers-howto:
// Audience must be a string with the following values:
// * App Engine: /projects/PROJECT_NUMBER/apps/PROJECT_ID
// * Compute Engine and Container Engine: /projects/PROJECT_NUMBER/global/backendServices/SERVICE_ID
type Audience string

// NewAudience returns an Audience from a string.
func NewAudience(u *string) *Audience {
	return (*Audience)(u)
}

// Sanitize normalizes the structure of the Audience's URL and validates it.
func (aud *Audience) Sanitize() error {
	return aud.Validate()
}

// Validate performs error checking on the Audience's URL.
func (aud *Audience) Validate() error {
	// TODO: Add actual validation
	return nil
}

// ParseAudience parses an Audience from a string.
func ParseAudience(rawAudience string) (*Audience, error) {
	aud := NewAudience(&rawAudience)
	if err := aud.Sanitize(); err != nil {
		return nil, err
	}
	return aud, nil
}
