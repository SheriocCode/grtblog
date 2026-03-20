package contract

import (
	"encoding/json"
	"strings"
	"time"
)

type CreateGalleryReq struct {
	Content     string     `json:"content" validate:"required"`
	Images      []string   `json:"images,omitempty"`
	IsPublished bool       `json:"isPublished"`
	IsTop       bool       `json:"isTop"`
	ExtInfo     *JSONRaw   `json:"extInfo,omitempty" swaggertype:"object"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}

type createGalleryReqJSON struct {
	Content     string   `json:"content"`
	Images      []string `json:"images"`
	IsPublished bool     `json:"isPublished"`
	IsTop       bool     `json:"isTop"`
	ExtInfo     *JSONRaw `json:"extInfo" swaggertype:"object"`
	CreatedAt   *string  `json:"createdAt"`
}

func (r *CreateGalleryReq) UnmarshalJSON(data []byte) error {
	var aux createGalleryReqJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	r.Content = aux.Content
	r.Images = aux.Images
	r.IsPublished = aux.IsPublished
	r.IsTop = aux.IsTop
	r.ExtInfo = aux.ExtInfo

	if aux.CreatedAt == nil {
		r.CreatedAt = nil
		return nil
	}
	if strings.TrimSpace(*aux.CreatedAt) == "" {
		now := time.Now()
		r.CreatedAt = &now
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, *aux.CreatedAt)
	if err != nil {
		return err
	}
	r.CreatedAt = &parsed
	return nil
}

type UpdateGalleryReq struct {
	Content     string     `json:"content" validate:"required"`
	Images      []string   `json:"images,omitempty"`
	IsPublished bool       `json:"isPublished"`
	IsTop       bool       `json:"isTop"`
	ExtInfo     *JSONRaw   `json:"extInfo,omitempty" swaggertype:"object"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}

type ListGalleriesReq struct {
	Page      int   `json:"page" validate:"min=1"`
	PageSize  int   `json:"pageSize" validate:"min=1,max=100"`
	Published *bool `json:"published,omitempty"`
}

type BatchSetGalleryPublishedReq struct {
	IDs         []int64 `json:"ids"`
	IsPublished bool    `json:"isPublished"`
}

type BatchSetGalleryTopReq struct {
	IDs   []int64 `json:"ids"`
	IsTop bool    `json:"isTop"`
}

type BatchDeleteGalleryReq struct {
	IDs []int64 `json:"ids"`
}