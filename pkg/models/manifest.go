package models

import (
	"time"
)

// Manifest defines spec for image manifest
type Manifest struct {
	Name         string    `json:"name"`
	Tag          string    `json:"tag"`
	Architecture string    `json:"architecture"`
	Layers       []Layer   `json:"fsLayers"`
	History      []History `json:"history"`
}

// History provides definition of container exec
type History struct {
	V1Compatibility string `json:"v1Compatibility"`
}

// V1Compatibility ..
type V1Compatibility struct {
	ID      string    `json:"id"`
	Cmd     string    `json:"Cmd"`
	Created time.Time `json:"created"`
}

// Layer defines representation of layer
type Layer struct {
	BlobSum string `json:"blobSum"`
}
