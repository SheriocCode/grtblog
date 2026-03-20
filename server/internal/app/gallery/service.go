package gallery

import (
	"context"
	"strings"
	"time"

	appEvent "github.com/grtsinry43/grtblog-v2/server/internal/app/event"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
)

type Service struct {
	repo   content.Repository
	events appEvent.Bus
}

func NewService(repo content.Repository, events appEvent.Bus) *Service {
	if events == nil {
		events = appEvent.NopBus{}
	}
	return &Service{repo: repo, events: events}
}

func (s *Service) CreateGallery(ctx context.Context, authorID int64, cmd CreateGalleryCmd) (*content.Gallery, error) {
	createdAt := time.Now()
	if cmd.CreatedAt != nil {
		createdAt = *cmd.CreatedAt
	}

	gallery := &content.Gallery{
		Content:     strings.TrimSpace(cmd.Content),
		ContentHash: content.GalleryContentHash(strings.TrimSpace(cmd.Content), normalizeImages(cmd.Images)),
		AuthorID:    authorID,
		Images:      normalizeImages(cmd.Images),
		IsPublished: cmd.IsPublished,
		IsTop:       cmd.IsTop,
		ExtInfo:     cmd.ExtInfo,
		CreatedAt:   createdAt,
	}

	if err := s.repo.CreateGallery(ctx, gallery); err != nil {
		return nil, err
	}

	now := time.Now()
	_ = s.events.Publish(ctx, GalleryCreated{ID: gallery.ID, AuthorID: gallery.AuthorID, Published: gallery.IsPublished, At: now})
	if gallery.IsPublished {
		_ = s.events.Publish(ctx, GalleryPublished{ID: gallery.ID, AuthorID: gallery.AuthorID, At: now})
	}

	return gallery, nil
}

func (s *Service) UpdateGallery(ctx context.Context, cmd UpdateGalleryCmd) (*content.Gallery, error) {
	existing, err := s.repo.GetGalleryByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	prevPublished := existing.IsPublished

	existing.Content = strings.TrimSpace(cmd.Content)
	existing.Images = normalizeImages(cmd.Images)
	existing.ContentHash = content.GalleryContentHash(existing.Content, existing.Images)
	existing.IsPublished = cmd.IsPublished
	existing.IsTop = cmd.IsTop
	existing.ExtInfo = cmd.ExtInfo
	if cmd.CreatedAt != nil {
		existing.CreatedAt = *cmd.CreatedAt
	}

	if err := s.repo.UpdateGallery(ctx, existing); err != nil {
		return nil, err
	}

	now := time.Now()
	_ = s.events.Publish(ctx, GalleryUpdated{ID: existing.ID, AuthorID: existing.AuthorID, Published: existing.IsPublished, ContentHash: existing.ContentHash, At: now})
	if !prevPublished && existing.IsPublished {
		_ = s.events.Publish(ctx, GalleryPublished{ID: existing.ID, AuthorID: existing.AuthorID, At: now})
	}
	if prevPublished && !existing.IsPublished {
		_ = s.events.Publish(ctx, GalleryUnpublished{ID: existing.ID, AuthorID: existing.AuthorID, At: now})
	}

	return existing, nil
}

func (s *Service) GetGalleryByID(ctx context.Context, id int64) (*content.Gallery, error) {
	return s.repo.GetGalleryByID(ctx, id)
}

func (s *Service) ListGalleries(ctx context.Context, options content.GalleryListOptionsInternal) ([]*content.Gallery, int64, error) {
	return s.repo.ListGalleries(ctx, options)
}

func (s *Service) ListPublicGalleries(ctx context.Context, options content.GalleryListOptions) ([]*content.Gallery, int64, error) {
	return s.repo.ListPublicGalleries(ctx, options)
}

func (s *Service) BatchSetPublished(ctx context.Context, cmd BatchSetPublishedCmd) error {
	ids := normalizeIDs(cmd.IDs)
	for _, id := range ids {
		item, err := s.repo.GetGalleryByID(ctx, id)
		if err != nil {
			return err
		}
		prevPublished := item.IsPublished
		if prevPublished == cmd.IsPublished {
			continue
		}
		item.IsPublished = cmd.IsPublished
		if err := s.repo.UpdateGallery(ctx, item); err != nil {
			return err
		}
		now := time.Now()
		_ = s.events.Publish(ctx, GalleryUpdated{ID: item.ID, AuthorID: item.AuthorID, Published: item.IsPublished, ContentHash: item.ContentHash, At: now})
		if cmd.IsPublished {
			_ = s.events.Publish(ctx, GalleryPublished{ID: item.ID, AuthorID: item.AuthorID, At: now})
		} else {
			_ = s.events.Publish(ctx, GalleryUnpublished{ID: item.ID, AuthorID: item.AuthorID, At: now})
		}
	}
	return nil
}

func (s *Service) BatchSetTop(ctx context.Context, cmd BatchSetTopCmd) error {
	ids := normalizeIDs(cmd.IDs)
	for _, id := range ids {
		item, err := s.repo.GetGalleryByID(ctx, id)
		if err != nil {
			return err
		}
		if item.IsTop == cmd.IsTop {
			continue
		}
		item.IsTop = cmd.IsTop
		if err := s.repo.UpdateGallery(ctx, item); err != nil {
			return err
		}
		_ = s.events.Publish(ctx, GalleryUpdated{ID: item.ID, AuthorID: item.AuthorID, Published: item.IsPublished, ContentHash: item.ContentHash, At: time.Now()})
	}
	return nil
}

func (s *Service) DeleteGallery(ctx context.Context, id int64) error {
	item, err := s.repo.GetGalleryByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.DeleteGallery(ctx, id); err != nil {
		return err
	}
	_ = s.events.Publish(ctx, GalleryDeleted{ID: item.ID, AuthorID: item.AuthorID, At: time.Now()})
	return nil
}

func (s *Service) BatchDelete(ctx context.Context, cmd BatchDeleteCmd) error {
	for _, id := range normalizeIDs(cmd.IDs) {
		if err := s.DeleteGallery(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func normalizeIDs(ids []int64) []int64 {
	seen := make(map[int64]struct{}, len(ids))
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func normalizeImages(images []string) []string {
	out := make([]string, 0, len(images))
	for _, image := range images {
		trimmed := strings.TrimSpace(image)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}