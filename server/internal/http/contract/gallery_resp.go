package contract

import "time"

type GalleryResp struct {
	ID          int64      `json:"id"`
	Content     string     `json:"content"`
	ContentHash string     `json:"contentHash"`
	AuthorID    int64      `json:"authorId"`
	Images      []string   `json:"images,omitempty"`
	ImageCount  int        `json:"imageCount"`
	IsPublished bool       `json:"isPublished"`
	IsTop       bool       `json:"isTop"`
	ExtInfo     *JSONRaw   `json:"extInfo,omitempty" swaggertype:"object"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type GalleryListItemResp struct {
	ID          int64     `json:"id"`
	Content     string    `json:"content"`
	ContentHash string    `json:"contentHash"`
	Images      []string  `json:"images,omitempty"`
	ImageCount  int       `json:"imageCount"`
	IsPublished bool      `json:"isPublished"`
	IsTop       bool      `json:"isTop"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type GalleryListResp struct {
	Items []GalleryListItemResp `json:"items"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}