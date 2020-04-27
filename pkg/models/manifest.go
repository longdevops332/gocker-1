package models

// Manifest defines spec for image manifest
type Manifest struct {
	Name         string   `json:"name"`
	Tag          string   `json:"tag"`
	Architecture string   `json:"architecture"`
	Layers       []Layer  `json:"fsLayers"`
	History      []string `json:"history"`
}

// Layer defines representation of layer
type Layer struct {
	BlobSum string `json:"blobSum"`
}
