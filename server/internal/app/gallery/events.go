package gallery

import "time"

type GalleryCreated struct {
	ID        int64
	AuthorID  int64
	Published bool
	At        time.Time
}

func (e GalleryCreated) Name() string { return "gallery.created" }
func (e GalleryCreated) OccurredAt() time.Time { return e.At }

type GalleryUpdated struct {
	ID          int64
	AuthorID    int64
	Published   bool
	ContentHash string
	At          time.Time
}

func (e GalleryUpdated) Name() string { return "gallery.updated" }
func (e GalleryUpdated) OccurredAt() time.Time { return e.At }

type GalleryPublished struct {
	ID       int64
	AuthorID int64
	At       time.Time
}

func (e GalleryPublished) Name() string { return "gallery.published" }
func (e GalleryPublished) OccurredAt() time.Time { return e.At }

type GalleryUnpublished struct {
	ID       int64
	AuthorID int64
	At       time.Time
}

func (e GalleryUnpublished) Name() string { return "gallery.unpublished" }
func (e GalleryUnpublished) OccurredAt() time.Time { return e.At }

type GalleryDeleted struct {
	ID       int64
	AuthorID int64
	At       time.Time
}

func (e GalleryDeleted) Name() string { return "gallery.deleted" }
func (e GalleryDeleted) OccurredAt() time.Time { return e.At }