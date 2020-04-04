package images

import (
	"github.com/pkg/errors"
)

// Pull provides pulling of the images
type Pull struct {
	tag string
}

// NewPull provides initialization on the pulling
func NewPull() *Pull {
	return &Pull{
		tag: "latest",
	}
}

// WithTag provides adding of tags for image
// its overrides `latest` tag of the image
func (p *Pull) WithTag(tag string) {
	p.tag = tag
}

// Do starts operation of pulling
func (p *Pull) Do() error {
	token, err := p.getToken()
	if err != nil {
		return errors.Wrap(err, "unable to get token")
	}
}

// getToken return token for auth
func (p *Pull) getToken() (string, error) {
	return "", nil
}
