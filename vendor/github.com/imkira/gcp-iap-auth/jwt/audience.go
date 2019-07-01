package jwt

import (
	"fmt"
	"strings"
)

// Audience is a string wrapper to provide validation logic for GCP IAP audience URLs.
// From the IAP docs at https://cloud.google.com/iap/docs/signed-headers-howto:
// Audience must be a string with the following values:
// * App Engine: /projects/PROJECT_NUMBER/apps/PROJECT_ID
// * Compute Engine and Container Engine: /projects/PROJECT_NUMBER/global/backendServices/SERVICE_ID
type Audience string

// NewAudience returns an Audience from a string.
func NewAudience(u string) *Audience {
	aud := Audience(u)
	return &aud
}

// Validate performs error checking on the Audience's URL.
func (aud *Audience) Validate() error {
	rawAudience := string(*aud)
	p := strings.SplitN(rawAudience, "/", 4)
	if len(p) != 4 {
		return fmt.Errorf("audience %q must follow the format \"/projects/PROJECT_NUMBER/SERVICE_DETAILS\"", rawAudience)
	}
	if p[0] != "" {
		return fmt.Errorf("audience %q should start with a slash", rawAudience)
	}
	if p[1] != "projects" {
		return fmt.Errorf("expecting \"projects\" but got %q in audience %q", p[1], rawAudience)
	}
	if projectNumber := p[2]; len(projectNumber) == 0 {
		return fmt.Errorf("audience %q must have a non-empty project number", rawAudience)
	}
	if serviceDetails := p[3]; len(serviceDetails) == 0 {
		return fmt.Errorf("audience %q is missing service details", rawAudience)
	}
	return nil
}

// ParseAudience parses an Audience from a string.
func ParseAudience(rawAudience string) (*Audience, error) {
	aud := NewAudience(rawAudience)
	if err := aud.Validate(); err != nil {
		return nil, err
	}
	return aud, nil
}
