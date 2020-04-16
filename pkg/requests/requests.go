package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// Get provides Get request to the server and returns response
func Get(url string, resp interface{}) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	return get(client, req, resp)
}

// GetWithAuth provides sending of Get request with auth token
func GetWithAuth(token, url string, resp interface{}) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "unable to get manifest")
	}
	defer res.Body.Close()
	bytesRead := 0
	buf := make([]byte, 1024*1024)
	for {
		n, err := res.Body.Read(buf)
		bytesRead += n

		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "errors reading response")
		}
		fmt.Println("NN: ", n)
	}
	return nil
}

// general method for send get request
func get(client *http.Client, r *http.Request, resp interface{}) error {
	r.Header.Set("Content-Type", "application/json")
	res, err := client.Do(r)
	if err != nil {
		return errors.Wrap(err, "unable to get manifest")
	}
	defer res.Body.Close()
	if err := decodeJSON(res.Body, resp); err != nil {
		return err
	}
	return nil
}

func decodeJSON(body io.ReadCloser, resp interface{}) error {
	err := json.NewDecoder(body).Decode(&resp)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal data")
	}
	return nil
}
