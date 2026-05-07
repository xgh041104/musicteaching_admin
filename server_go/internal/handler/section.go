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

type SectionHandler interface {
	EditerUploadFile(ctx *gin.Context)
	AddSection(ctx *gin.Context)
	DelSection(ctx *gin.Context)
	EditSection(ctx *gin.Context)
	QuerySectionBySectionId(ctx *gin.Context)
	QuerySectionByChapterId(ctx *gin.Context)
}

type sectionHandler struct {
	*Handler
	sectionService service.SectionService
}

func NewSectionHandler(Handler *Handler, sectionService service.SectionService) SectionHandler {
	return &sectionHandler{
		Handler:        Handler,
		sectionService: sectionService,
	}
}

// 上传富文本编辑器图片
func (h *sectionHandler) EditerUploadFile(ctx *gin.Context) {
	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, err := h.sectionService.EditerUploadFile(files)
	//记录获取的学校
	h.logger.Info("EditerUploadFile", zap.Any("EditerUploadFile", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, data)
}

// 添加小节
func (h *sectionHandler) AddSection(ctx *gin.Context) {

	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var sectionparam model.SectionEdit

	json.Unmarshal([]byte(data), &sectionparam)

	err := h.sectionService.AddSection(sectionparam, files)
	//添加小节
	h.logger.Info("AddSection", zap.Any("AddSection", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, nil)
}

// 删除小节
func (h *sectionHandler) DelSection(ctx *gin.Context) {

	var sectionparam model.Section

	if err := ctx.ShouldBind(&sectionparam); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.sectionService.DelSection(sectionparam)
	//删除小节
	h.logger.Info("AddSection", zap.Any("AddSection", sectionparam))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, nil)
}

// 修改小节
func (h *sectionHandler) EditSection(ctx *gin.Context) {
	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var sectionparam model.SectionEdit

	json.Unmarshal([]byte(data), &sectionparam)

	err := h.sectionService.EditSection(sectionparam, files)
	//添加小节
	h.logger.Info("EditSection", zap.Any("EditSection", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, nil)
}

func (h *sectionHandler) QuerySectionBySectionId(ctx *gin.Context) {
	var params struct {
		SectionId int `form:"sectionId"` // 请求参数 courseId
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.sectionService.QuerySectionBySectionId(params.SectionId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), data)
		return
	}
	resp.HandleSuccess(ctx, data)
}

func (h *sectionHandler) QuerySectionByChapterId(ctx *gin.Context) {
	var params struct {
		ChapterId int `form:"chapterId"` // 请求参数 courseId
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.sectionService.QuerySectionByChapterId(params.ChapterId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), data)
		return
	}
	resp.HandleSuccess(ctx, data)
}
