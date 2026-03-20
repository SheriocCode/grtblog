package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/gallery"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
)

func registerGalleryPublicRoutes(v2 fiber.Router, deps Dependencies) {
	galleryHandler := newGalleryHandler(deps)
	publicGroup := v2.Group("/galleries")
	publicGroup.Get("/", galleryHandler.ListGalleries)
}

func registerGalleryAuthRoutes(v2 fiber.Router, deps Dependencies) {
	galleryHandler := newGalleryHandler(deps)
	identityRepo := persistence.NewIdentityRepository(deps.DB)
	adminTokenRepo := persistence.NewAdminTokenRepository(deps.DB)

	authGroup := v2.Group("/galleries", middleware.RequireAuth(deps.JWTManager, identityRepo, adminTokenRepo), middleware.RequireAdmin(identityRepo))
	authGroup.Post("/", galleryHandler.CreateGallery)
	authGroup.Put("/:id", galleryHandler.UpdateGallery)
	authGroup.Delete("/:id", galleryHandler.DeleteGallery)

	adminGroup := v2.Group("/admin", middleware.RequireAuth(deps.JWTManager, identityRepo, adminTokenRepo), middleware.RequireAdmin(identityRepo))
	adminGroup.Get("/galleries/:id", galleryHandler.GetGalleryAdmin)
	adminGroup.Get("/galleries", galleryHandler.ListGalleriesAdmin)
	adminGroup.Put("/galleries/published", galleryHandler.BatchSetGalleryPublished)
	adminGroup.Put("/galleries/top", galleryHandler.BatchSetGalleryTop)
	adminGroup.Post("/galleries/batch-delete", galleryHandler.BatchDeleteGalleries)
}

func newGalleryHandler(deps Dependencies) *handler.GalleryHandler {
	contentRepo := persistence.NewContentRepository(deps.DB)
	gallerySvc := gallery.NewService(contentRepo, deps.EventBus)
	return handler.NewGalleryHandler(gallerySvc)
}