package network

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Get provides Get request to teh server and returns response
func Get(url string, resp interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "unable to get auth")
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal data")
	}
	return nil
}
