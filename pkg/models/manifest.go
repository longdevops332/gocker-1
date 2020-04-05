package models

// Manifest defines spec for image manifest
type Manifest struct {
	Name   string  `json:"manifest"`
	Layers []Layer `json:"fsLayers"`
}

// Layer defines representation of layer
type Layer struct {
	BlobSum string `json:"blobSum"`
}
