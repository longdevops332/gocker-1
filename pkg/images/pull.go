package images

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/saromanov/gocker/pkg/models"
	"github.com/saromanov/gocker/pkg/requests"
)

// Pull provides pulling of the images
type Pull struct {
	tag   string
	image string
}

// NewPull provides initialization on the pulling
func NewPull(image string) *Pull {
	return &Pull{
		image: image,
		tag:   "latest",
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
	fmt.Println("TOKEN: ", token)
	return nil
}

// getToken return token for auth
func (p *Pull) getToken() (string, error) {
	var t *models.Auth
	err := requests.Get(fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s/%s:pull", "name", p.image), &t)
	if err != nil {
		return "", errors.Wrap(err, "unable to get auth")
	}
	if t == nil {
		return "", errors.New("unable to unmarshal token")
	}
	return t.Token, nil
}
