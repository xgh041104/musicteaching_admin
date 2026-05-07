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

type GiveScoreCategoryHandler interface {
	AddGiveScoreCategory(ctx *gin.Context)
	EditGiveScoreCategory(ctx *gin.Context)
	DelGiveScoreCategory(ctx *gin.Context)
	QueryGiveScoreCategory(ctx *gin.Context)
	QueryGiveScoreCategoryTree(ctx *gin.Context)
	QueryGiveScoreCategoryParentNodeByParentId(ctx *gin.Context)
	QueryGiveScoreCategoryChildNodesById(ctx *gin.Context)
}

type giveScoreCategoryHandler struct {
	*Handler
	giveScoreCategoryService service.GiveScoreCategoryService
}

func NewgiveScoreCategoryHandler(handler *Handler, giveScoreCategoryService service.GiveScoreCategoryService) GiveScoreCategoryHandler {
	return &giveScoreCategoryHandler{
		Handler:                  handler,
		giveScoreCategoryService: giveScoreCategoryService,
	}
}

func (h *giveScoreCategoryHandler) AddGiveScoreCategory(ctx *gin.Context) {

	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var giveScoreCategoryparam model.GiveScoreCategory

	json.Unmarshal([]byte(data), &giveScoreCategoryparam)

	err := h.giveScoreCategoryService.AddGiveScoreCategory(giveScoreCategoryparam, files)

	h.logger.Info("AddGiveScoreCategory", zap.Any("AddGiveScoreCategory", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, nil)

}

func (h *giveScoreCategoryHandler) EditGiveScoreCategory(ctx *gin.Context) {

	files, _ := ctx.MultipartForm()

	data, _ := ctx.GetPostForm("data")

	var giveScoreCategoryparam model.GiveScoreCategory

	json.Unmarshal([]byte(data), &giveScoreCategoryparam)

	err := h.giveScoreCategoryService.EditGiveScoreCategory(giveScoreCategoryparam, files)

	h.logger.Info("EditGiveScoreCategory", zap.Any("EditGiveScoreCategory", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, nil)
}

func (h *giveScoreCategoryHandler) DelGiveScoreCategory(ctx *gin.Context) {
	var params struct {
		GiveScoreCategoryId int64 `json:"giveScoreCategoryId" binding:"required"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.giveScoreCategoryService.DelGiveScoreCategory(params.GiveScoreCategoryId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}

func (h *giveScoreCategoryHandler) QueryGiveScoreCategory(ctx *gin.Context) {
	var params struct {
		SchoolId int64 `form:"schoolId"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.giveScoreCategoryService.QueryGiveScoreCategory(params.SchoolId)

	h.logger.Info("QueryGiveScoreCategory", zap.Any("QueryGiveScoreCategory", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, data)
}

func (h *giveScoreCategoryHandler) QueryGiveScoreCategoryTree(ctx *gin.Context) {
	var params struct {
		SchoolId int64 `form:"schoolId"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.giveScoreCategoryService.QueryGiveScoreCategoryTree(params.SchoolId)

	h.logger.Info("QueryGiveScoreCategoryTree", zap.Any("QueryGiveScoreCategoryTree", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), data)
		return
	}

	resp.HandleSuccess(ctx, data)
}

func (h *giveScoreCategoryHandler) QueryGiveScoreCategoryParentNodeByParentId(ctx *gin.Context) {
	var params struct {
		GiveScoreCategoryParentId int64 `form:"giveScoreCategoryParentId" binding:"required"`
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.giveScoreCategoryService.QueryGiveScoreCategoryParentNodeByParentId(params.GiveScoreCategoryParentId)

	h.logger.Info("QueryGiveScoreCategoryParentNodeByParentId", zap.Any("QueryGiveScoreCategoryParentNodeByParentId", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), data)
		return
	}

	resp.HandleSuccess(ctx, data)
}

func (h *giveScoreCategoryHandler) QueryGiveScoreCategoryChildNodesById(ctx *gin.Context) {
	var params struct {
		GiveScoreCategoryId int64 `form:"giveScoreCategoryId" binding:"required"`
		SchoolId            int64 `form:"schoolId"`
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.giveScoreCategoryService.QueryGiveScoreCategoryChildNodesById(params.GiveScoreCategoryId, params.SchoolId)

	h.logger.Info("QueryGiveScoreCategoryChildNodesById", zap.Any("QueryGiveScoreCategoryChildNodesById", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), data)
		return
	}

	resp.HandleSuccess(ctx, data)
}
