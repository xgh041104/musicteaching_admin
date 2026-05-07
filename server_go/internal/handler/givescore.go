package handler

import (
	"encoding/json"
	"net/http"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/pkg/helper/resp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GiveScoreHandler interface {
	QueryGiveScore(ctx *gin.Context)
	QueryGiveTeacherScoreByGiveScoreThemeId(ctx *gin.Context)
	QueryGiveOneselfScoreByGiveScoreThemeId(ctx *gin.Context)
	QueryGiveMutualScoreByGiveScoreThemeId(ctx *gin.Context)
	AddGiveScore(ctx *gin.Context)
	EditGiveScore(ctx *gin.Context)
	DelGiveSCore(ctx *gin.Context)
}

type giveScoreHandler struct {
	*Handler
	giveScoreService service.GiveScoreService
}

func NewgiveScoreHandler(handler *Handler, giveScoreService service.GiveScoreService) GiveScoreHandler {
	return &giveScoreHandler{
		Handler:          handler,
		giveScoreService: giveScoreService,
	}
}

func (h *giveScoreHandler) QueryGiveScore(ctx *gin.Context) {
	var params struct {
		GiveScoreThemeId int64 `form:"giveScoreThemeId"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.giveScoreService.QueryGiveScore(params.GiveScoreThemeId)

	h.logger.Info("QueryGiveScore ", zap.Any("QueryGiverScore", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, data)
}
func (h *giveScoreHandler) QueryGiveTeacherScoreByGiveScoreThemeId(ctx *gin.Context) {
	var params struct {
		GiveScoreThemeId int64 `form:"giveScoreThemeId" binding:"required"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.giveScoreService.QueryGiveTeacherScoreByGiveScoreThemeId(params.GiveScoreThemeId)

	h.logger.Info("QueryGiveTeacherScore ", zap.Any("QueryGiveTeacherScore", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, data)
}

func (h *giveScoreHandler) QueryGiveOneselfScoreByGiveScoreThemeId(ctx *gin.Context) {
	var params struct {
		GiveScoreThemeId int64 `form:"giveScoreThemeId" binding:"required"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.giveScoreService.QueryGiveOneselfScoreByGiveScoreThemeId(params.GiveScoreThemeId)

	h.logger.Info("QueryGiveOneselfScore ", zap.Any("QueryGiveOneselfScore", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, data)
}

func (h *giveScoreHandler) QueryGiveMutualScoreByGiveScoreThemeId(ctx *gin.Context) {
	var params struct {
		GiveScoreThemeId int64 `form:"giveScoreThemeId" binding:"required"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.giveScoreService.QueryGiveMutualScoreByGiveScoreThemeId(params.GiveScoreThemeId)

	h.logger.Info("QueryGiveMutualScore ", zap.Any("QueryGiveMutualScore", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, data)
}

func (h *giveScoreHandler) AddGiveScore(ctx *gin.Context) {
	files, _ := ctx.MultipartForm()

	data, _ := ctx.GetPostForm("data")

	var giveScoreparam model.GiveScore

	json.Unmarshal([]byte(data), &giveScoreparam)

	h.logger.Info("QueryGiveMutualScore ", zap.Any("QueryGiveMutualScore", data))
	err := h.giveScoreService.AddGiveScore(giveScoreparam, files)

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)

}

func (h *giveScoreHandler) EditGiveScore(ctx *gin.Context) {
	var giveScore model.GiveScore

	if err := ctx.ShouldBind(&giveScore); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.giveScoreService.EditGiveScore(giveScore)

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}

func (h *giveScoreHandler) DelGiveSCore(ctx *gin.Context) {
	var params struct {
		GiveScoreId int64 `json:"giveScoreId" binding:"required"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.giveScoreService.DelGiveSCore(params.GiveScoreId)

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}
