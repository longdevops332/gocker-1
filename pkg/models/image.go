package models

// Image defines basic image representation
type Image struct {
	ID      int64
	Name    string
	Version string
	Size    uint64
	Path    string
}
