package gallery

import "time"

type CreateGalleryCmd struct {
	Content     string
	Images      []string
	IsPublished bool
	IsTop       bool
	ExtInfo     []byte
	CreatedAt   *time.Time
}

type UpdateGalleryCmd struct {
	ID          int64
	Content     string
	Images      []string
	IsPublished bool
	IsTop       bool
	ExtInfo     []byte
	CreatedAt   *time.Time
}

type BatchSetPublishedCmd struct {
	IDs         []int64
	IsPublished bool
}

type BatchSetTopCmd struct {
	IDs   []int64
	IsTop bool
}

type BatchDeleteCmd struct {
	IDs []int64
}