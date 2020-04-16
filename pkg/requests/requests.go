package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Get provides Get request to teh server and returns response
func Get(url string, resp interface{}) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	return get(client, req, resp)
}

// GetWithAuth provides sending of Get request with auth token
func GetWithAuth(token, url string, resp interface{}) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return get(client, req, resp)
}

// general method for send get request
func get(client *http.Client, r *http.Request, resp interface{}) error {
	res, err := client.Do(r)
	if err != nil {
		return errors.Wrap(err, "unable to get manifest")
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal data")
	}
	return nil
}
