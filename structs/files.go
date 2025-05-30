package structs

import "time"

type FileEntry struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	IsDir        bool      `json:"isDir"`
	Size         int64     `json:"size,omitempty"`
	LastModified time.Time `json:"lastModified,omitzero"`
}

type CreateItemRequest struct {
	Path    string `json:"path"`
	IsDir   bool   `json:"isDir"`
	Content string `json:"content,omitempty"`
}
