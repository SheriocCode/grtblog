package handler

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/gallery"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/middleware"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type GalleryHandler struct {
	svc *gallery.Service
}

func NewGalleryHandler(svc *gallery.Service) *GalleryHandler {
	return &GalleryHandler{svc: svc}
}

func (h *GalleryHandler) CreateGallery(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	var req contract.CreateGalleryReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if strings.TrimSpace(req.Content) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "内容不能为空")
	}
	extInfo, err := parseExtInfo(req.ExtInfo)
	if err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "extInfo格式错误", err)
	}

	createdGallery, err := h.svc.CreateGallery(c.Context(), claims.UserID, gallery.CreateGalleryCmd{
		Content:     req.Content,
		Images:      req.Images,
		IsPublished: req.IsPublished,
		IsTop:       req.IsTop,
		ExtInfo:     extInfo,
		CreatedAt:   req.CreatedAt,
	})
	if err != nil {
		return err
	}

	Audit(c, "gallery.create", map[string]any{"galleryId": createdGallery.ID, "userId": claims.UserID})
	return response.SuccessWithMessage(c, toGalleryResp(createdGallery), "日常创建成功")
}

func (h *GalleryHandler) UpdateGallery(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的日常ID")
	}

	var req contract.UpdateGalleryReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if strings.TrimSpace(req.Content) == "" {
		return response.NewBizErrorWithMsg(response.ParamsError, "内容不能为空")
	}
	extInfo, err := parseExtInfo(req.ExtInfo)
	if err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "extInfo格式错误", err)
	}

	updatedGallery, err := h.svc.UpdateGallery(c.Context(), gallery.UpdateGalleryCmd{
		ID:          id,
		Content:     req.Content,
		Images:      req.Images,
		IsPublished: req.IsPublished,
		IsTop:       req.IsTop,
		ExtInfo:     extInfo,
		CreatedAt:   req.CreatedAt,
	})
	if err != nil {
		if errors.Is(err, content.ErrGalleryNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "日常不存在")
		}
		return err
	}

	Audit(c, "gallery.update", map[string]any{"galleryId": updatedGallery.ID, "userId": claims.UserID})
	return response.SuccessWithMessage(c, toGalleryResp(updatedGallery), "日常更新成功")
}

func (h *GalleryHandler) GetGalleryAdmin(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的日常ID")
	}
	item, err := h.svc.GetGalleryByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, content.ErrGalleryNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "日常不存在")
		}
		return err
	}
	return response.Success(c, toGalleryResp(item))
}

func (h *GalleryHandler) ListGalleries(c *fiber.Ctx) error {
	query := buildGalleryListQuery(c)
	items, total, err := h.svc.ListPublicGalleries(c.Context(), content.GalleryListOptions{Page: query.Page, PageSize: query.PageSize})
	if err != nil {
		return err
	}
	return response.Success(c, toGalleryListResp(items, total, query.Page, query.PageSize))
}

func (h *GalleryHandler) ListGalleriesAdmin(c *fiber.Ctx) error {
	query := buildGalleryListQuery(c)
	items, total, err := h.svc.ListGalleries(c.Context(), content.GalleryListOptionsInternal{Page: query.Page, PageSize: query.PageSize, Published: query.Published})
	if err != nil {
		return err
	}
	return response.Success(c, toGalleryListResp(items, total, query.Page, query.PageSize))
}

func (h *GalleryHandler) ListRecentPublicGalleries(c *fiber.Ctx) error {
	const page = 1
	const size = 6
	items, total, err := h.svc.ListPublicGalleries(c.Context(), content.GalleryListOptions{Page: page, PageSize: size})
	if err != nil {
		return err
	}
	return response.Success(c, toGalleryListResp(items, total, page, size))
}

func (h *GalleryHandler) DeleteGallery(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ParamsError, "无效的日常ID")
	}
	if err := h.svc.DeleteGallery(c.Context(), id); err != nil {
		if errors.Is(err, content.ErrGalleryNotFound) {
			return response.NewBizErrorWithMsg(response.NotFound, "日常不存在")
		}
		return err
	}
	Audit(c, "gallery.delete", map[string]any{"galleryId": id, "userId": claims.UserID})
	return response.SuccessWithMessage[any](c, nil, "日常删除成功")
}

func (h *GalleryHandler) BatchSetGalleryPublished(c *fiber.Ctx) error {
	var req contract.BatchSetGalleryPublishedReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if len(req.IDs) == 0 {
		return response.NewBizErrorWithMsg(response.ParamsError, "ids 不能为空")
	}
	if err := h.svc.BatchSetPublished(c.Context(), gallery.BatchSetPublishedCmd{IDs: req.IDs, IsPublished: req.IsPublished}); err != nil {
		return err
	}
	if req.IsPublished {
		return response.SuccessWithMessage[any](c, nil, "日常发布状态已批量更新为已发布")
	}
	return response.SuccessWithMessage[any](c, nil, "日常发布状态已批量更新为未发布")
}

func (h *GalleryHandler) BatchSetGalleryTop(c *fiber.Ctx) error {
	var req contract.BatchSetGalleryTopReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if len(req.IDs) == 0 {
		return response.NewBizErrorWithMsg(response.ParamsError, "ids 不能为空")
	}
	if err := h.svc.BatchSetTop(c.Context(), gallery.BatchSetTopCmd{IDs: req.IDs, IsTop: req.IsTop}); err != nil {
		return err
	}
	if req.IsTop {
		return response.SuccessWithMessage[any](c, nil, "日常置顶状态已批量更新为置顶")
	}
	return response.SuccessWithMessage[any](c, nil, "日常置顶状态已批量更新为取消置顶")
}

func (h *GalleryHandler) BatchDeleteGalleries(c *fiber.Ctx) error {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return response.ErrorFromBiz[any](c, response.NotLogin)
	}
	var req contract.BatchDeleteGalleryReq
	if err := c.BodyParser(&req); err != nil {
		return response.NewBizErrorWithCause(response.ParamsError, "请求体解析失败", err)
	}
	if len(req.IDs) == 0 {
		return response.NewBizErrorWithMsg(response.ParamsError, "ids 不能为空")
	}
	if err := h.svc.BatchDelete(c.Context(), gallery.BatchDeleteCmd{IDs: req.IDs}); err != nil {
		return err
	}
	Audit(c, "gallery.batch_delete", map[string]any{"galleryIds": req.IDs, "userId": claims.UserID})
	return response.SuccessWithMessage[any](c, nil, "日常批量删除成功")
}

func buildGalleryListQuery(c *fiber.Ctx) contract.ListGalleriesReq {
	query := contract.ListGalleriesReq{Page: 1, PageSize: 12}
	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		query.Page = page
	}
	if pageSize, err := strconv.Atoi(c.Query("pageSize", "12")); err == nil && pageSize > 0 && pageSize <= 100 {
		query.PageSize = pageSize
	}
	if publishedStr := c.Query("published"); publishedStr != "" {
		if published, err := strconv.ParseBool(publishedStr); err == nil {
			query.Published = &published
		}
	}
	return query
}

func toGalleryResp(item *content.Gallery) contract.GalleryResp {
	images := append([]string(nil), item.Images...)
	return contract.GalleryResp{
		ID:          item.ID,
		Content:     item.Content,
		ContentHash: item.ContentHash,
		AuthorID:    item.AuthorID,
		Images:      images,
		ImageCount:  len(images),
		IsPublished: item.IsPublished,
		IsTop:       item.IsTop,
		ExtInfo:     jsonRawFromBytes(item.ExtInfo),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func toGalleryListResp(items []*content.Gallery, total int64, page int, pageSize int) contract.GalleryListResp {
	respItems := make([]contract.GalleryListItemResp, len(items))
	for i, item := range items {
		images := append([]string(nil), item.Images...)
		respItems[i] = contract.GalleryListItemResp{
			ID:          item.ID,
			Content:     item.Content,
			ContentHash: item.ContentHash,
			Images:      images,
			ImageCount:  len(images),
			IsPublished: item.IsPublished,
			IsTop:       item.IsTop,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}
	return contract.GalleryListResp{Items: respItems, Total: total, Page: page, Size: pageSize}
}