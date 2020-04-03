package models

// Image defines basic image representation
type Image struct {
	ID      int64 `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Size    uint64 `json:"size"`
	Path    string `json:"path"`
}
