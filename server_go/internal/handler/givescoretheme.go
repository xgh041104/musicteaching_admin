package handler

import (
	"net/http"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/pkg/helper/resp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GiveScoreThemeHandler interface {
	QueryGiveScoreThemeByGiveScoreCategoryId(ctx *gin.Context)
	QueryGiveScoreTheme(ctx *gin.Context)
	AddGiveScoreTheme(ctx *gin.Context)
	EditGiveScoreTheme(ctx *gin.Context)
	DelGiveScoreTheme(ctx *gin.Context)
}

type giveScoreThemeHandler struct {
	*Handler
	giveScoreThemeService service.GiveScoreThemeService
}

func NewgiveScoreThemeHandler(handler *Handler, giveScoreThemeService service.GiveScoreThemeService) GiveScoreThemeHandler {
	return &giveScoreThemeHandler{
		Handler:               handler,
		giveScoreThemeService: giveScoreThemeService,
	}
}

func (h *giveScoreThemeHandler) AddGiveScoreTheme(ctx *gin.Context) {
	var giveScoreTheme model.GiveScoreTheme
	if err := ctx.ShouldBind(&giveScoreTheme); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.giveScoreThemeService.AddGiveScoreTheme(giveScoreTheme)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}

func (h *giveScoreThemeHandler) EditGiveScoreTheme(ctx *gin.Context) {
	var giveScoreTheme model.GiveScoreTheme
	if err := ctx.ShouldBind(&giveScoreTheme); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.giveScoreThemeService.EditGiveScoreTheme(giveScoreTheme)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}

func (h *giveScoreThemeHandler) DelGiveScoreTheme(ctx *gin.Context) {
	var params struct {
		GiveScoreThemeId int64 `form:"giveScoreThemeId" binding:"required"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.giveScoreThemeService.DelGiveScoreTheme(params.GiveScoreThemeId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}

func (h *giveScoreThemeHandler) QueryGiveScoreThemeByGiveScoreCategoryId(ctx *gin.Context) {
	var params struct {
		GiveScoreCategoryId int64 `form:"giveScoreCategoryId" binding:"required"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.giveScoreThemeService.QueryGiveScoreThemeByGiveScoreCategoryId(params.GiveScoreCategoryId)

	h.logger.Info("QueryGiveScoreTheme", zap.Any("QueryGiveScoreTheme", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, data)
}

func (h *giveScoreThemeHandler) QueryGiveScoreTheme(ctx *gin.Context) {
	var params struct {
		SchoolId int64 `form:"schoolId"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.giveScoreThemeService.QueryGiveScoreTheme(params.SchoolId)

	h.logger.Info("QueryGiveScoreTheme", zap.Any("QueryGiveScoreTheme", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, data)
}
