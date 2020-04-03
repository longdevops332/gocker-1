package images

import (
	"encoding/json"
	"strings"
	"fmt"
	"io/ioutil"
	"github.com/saromanov/gocker/pkg/models"
	"github.com/pkg/errors"
)

// List returns list of images
func List(dir string) ([]models.Image, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to read dir: %s", dir))
	}
	images := []models.Image{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		img, err := prepareImage(f.Name())
		if err != nil {
			return nil, errors.Wrap(err, "unable to get prepared image")
		}
		img.Path = f.Name()

	}
	return images, nil
}

func prepareImage(path string)(models.Image, error) {
	var img models.Image
	f, err := ioutil.ReadAll(path)
	if err != nil {
		return img, errors.Wrap(err, "unable to read image")
	}
	if err := json.Unmarshal(f, &img); err != nil {
		return img, errors,.Wrap(err, "unable to unmarshal image")
	}
	img.Size = sizeOfFmt(1000)
	return img, nil
}
