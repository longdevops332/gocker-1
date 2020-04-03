package images

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

// List returns list of images
func List(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to read dir: %s", dir))
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
	}
}
