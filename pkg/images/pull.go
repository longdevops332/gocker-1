package images

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	tag           string
	image         string
	library       string
	baseDirectory string
}

// NewPull provides initialization on the pulling
func NewPull(img string) *Pull {
	library, image := splitInputImage(img)
	baseDir, err := createBaseDirectory()
	if err != nil {
		panic(err)
	}
	return &Pull{
		image:         image,
		tag:           "latest",
		library:       library,
		baseDirectory: baseDir,
	}
}

func splitInputImage(img string) (string, string) {
	result := strings.Split(img, "/")
	if len(result) == 1 {
		return "library", result[0]
	}
	return result[0], result[1]
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

	manifest, err := p.getManifest(token, p.library, p.image, p.tag)
	if err != nil {
		return errors.Wrap(err, "unable to get manigest data")
	}
	if err := preparePulling(p.baseDirectory, manifest); err != nil {
		return errors.Wrap(err, "failed to write manifest file")
	}
	signs := getLayerSigns(manifest)
	if len(signs) == 0 {
		return errors.New("layers signs is not defined")
	}
	for sig := range signs {
		fmt.Printf("fetching the layer: %s\n", sig)
		url := fmt.Sprintf("%s/%s/%s/blobs/%s", registryURL, p.library, p.image, sig)
		var resp map[string]interface{}
		filePath := path.Join(p.baseDirectory, "layers")
		if err := requests.StreamToFile(filePath, token, url, &resp); err != nil {
			return errors.Wrap(err, "unable to get content")
		}
		fmt.Println(resp)
	}
	return nil
}

// getToken return token for auth
func (p *Pull) getToken() (string, error) {
	var t *models.Auth
	url := fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s/%s:pull", p.library, p.image)
	err := requests.Get(url, &t)
	if err != nil {
		return "", errors.Wrap(err, "unable to get auth")
	}
	if t == nil {
		return "", errors.New("unable to unmarshal token")
	}
	return t.Token, nil
}

// getManifest returns manifest of the image
func (p *Pull) getManifest(token, library, image, tag string) (*models.Manifest, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/%s/%s/manifests/%s", registryURL, library, image, tag)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get manifest")
	}
	var m *models.Manifest
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return nil, errors.Wrap(err, "unable to decode response")
	}
	if m == nil || len(m.Layers) == 0 {
		return nil, errors.New("manifest file is empty")
	}
	return m, nil
}

// preparePulling provides writing of manifest to the file
// also, its creating supported directories if this is not exists
func preparePulling(baseDir string, m *models.Manifest) error {
	if m.Name == "" {
		return errors.New("name of the image is not defined")
	}
	imageName := strings.Replace(m.Name, "/", "_", -1)
	pathDir := path.Join(baseDir, fmt.Sprintf("%s.json", imageName))
	data, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "unable to marshal manifest")
	}
	if err := ioutil.WriteFile(pathDir, data, 0664); err != nil {
		return errors.Wrap(err, "unable to write to file")
	}

	if layerPath, err := createSubDir(baseDir, imageName, "layers"); err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to create layer dir: %s", layerPath))
	}
	contentsPath, err := createSubDir(baseDir, imageName, "layers/contents")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to create content dir: %s", contentsPath))
	}

	return nil
}

// createSunDir provides creating of directory for image layers
func createSubDir(basePath, image, subDir string) (string, error) {
	layersPath := path.Join(basePath, image)
	layersPath = path.Join(layersPath, subDir)
	if _, err := os.Stat(layersPath); !os.IsNotExist(err) {
		return layersPath, nil
	}
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

// createBaseDirectory returns main dir for store images
// or returns default one
// its create base dir if this is not exists
func createBaseDirectory() (string, error) {
	baseDir := os.Getenv("GOCKER_BASE_DIR")
	if baseDir == "" {
		baseDir = "gocker-images"
	}
	if _, err := os.Stat(baseDir); !os.IsNotExist(err) {
		return baseDir, nil
	}
	if err := os.Mkdir(baseDir, os.ModePerm); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("unable to create base dir: %s", baseDir))
	}
	return baseDir, nil
}
