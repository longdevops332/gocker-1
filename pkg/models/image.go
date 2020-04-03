package models

// Image defines basic image representation
type Image struct {
	ID      int64 `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Size    string `json:"size"`
	Path    string `json:"path"`
}
