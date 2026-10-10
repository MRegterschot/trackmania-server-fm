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

type MoveItem struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type MoveItemsRequest struct {
	Items []MoveItem `json:"items"`
}
