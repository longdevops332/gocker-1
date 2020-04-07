package images

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/saromanov/gocker/pkg/models"
	"github.com/saromanov/gocker/pkg/requests"
)

const registryURL = "https://registry-1.docker.io/v2"

// Pull provides pulling of the images
type Pull struct {
	tag     string
	image   string
	library string
}

// NewPull provides initialization on the pulling
func NewPull(image, library string) *Pull {
	if library == "" {
		library = "lib"
	}
	return &Pull{
		image:   image,
		tag:     "latest",
		library: library,
	}
}

// WithTag provides adding of tags for image
// its overrides `latest` tag of the image
func (p *Pull) WithTag(tag string) {
	p.tag = tag
}

// Do starts operation of pulling
// Its pull image by layer and store it on tar file
func (p *Pull) Do() error {
	token, err := p.getToken()
	if err != nil {
		return errors.Wrap(err, "unable to get token")
	}
	fmt.Println("TOKEN: ", token)

	manifest, err := p.getManifest(p.library, p.image, p.tag)
	if err != nil {
		return errors.Wrap(err, "unable to get manigest data")
	}
	if err := writeManifestFile(manifest); err != nil {
		return errors.Wrap(err, "failed to write manifest file")
	}
	if layerPath, err := createSubDir("data", p.image, "layers"); err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to create layer dir: %s", layerPath))
	}
	contentsPath, err := createSubDir("data", p.image, "contents")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to create content dir: %s", contentsPath))
	}
	signs := getLayerSigns(manifest)
	for sig := range signs {
		url := fmt.Sprintf("%s/%s/%s/blobs/%s", registryURL, p.library, p.image, p.tag, sig)
		var resp map[string]interface{}
		if err := requests.Get(url, &resp); err != nil {
			return errors.Wrap(err, "unable to get content")
		}
	}
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

// getManifest returns manifest of the image
func (p *Pull) getManifest(library, image, tag string) (*models.Manifest, error) {
	var m *models.Manifest
	err := requests.Get(fmt.Sprintf("%s/%s/%s/manifests/%s", registryURL, library, image, tag), &m)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get manifest")
	}
	return m, nil
}

func writeManifestFile(m *models.Manifest) error {
	imageName := strings.Replace(m.Name, "/", "_", -1)
	data, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "unable to marshal manifest")
	}
	if err := ioutil.WriteFile(imageName, data, 0664); err != nil {
		return errors.Wrap(err, "unable to write to file")
	}
	return nil
}

// createSunDir provides creating of directory for image layers
func createSubDir(basePath, image, subDir string) (string, error) {
	layersPath := path.Join(basePath, image)
	layersPath = path.Join(layersPath, subDir)
	return layersPath, os.MkdirAll(layersPath, os.ModePerm)
}

// getLayerSigns returns signatures of layers
func getLayerSigns(m *models.Manifest) map[string]bool {
	result := make(map[string]bool)
	for _, l := range m.Layers {
		result[l.BlobSum] = true
	}
	return result
}
