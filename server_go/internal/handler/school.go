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

type SchoolHandler interface {
	AddSchool(ctx *gin.Context)
	EditSchool(ctx *gin.Context)
	DelSchool(ctx *gin.Context)
	QuerySchoolAll(ctx *gin.Context)
}

type schoolHandler struct {
	*Handler
	schoolService service.SchoolService
}

func NewSchoolHandler(handler *Handler, schoolService service.SchoolService) SchoolHandler {
	return &schoolHandler{
		Handler:       handler,
		schoolService: schoolService,
	}
}

func (h *schoolHandler) AddSchool(ctx *gin.Context) {

	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var schoolparam model.School

	json.Unmarshal([]byte(data), &schoolparam)

	err := h.schoolService.AddSchool(schoolparam, files)
	//记录获取的学校
	h.logger.Info("AddSchool", zap.Any("AddSchool", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, nil)
}

// 修改学校
func (h *schoolHandler) EditSchool(ctx *gin.Context) {
	//定义绑定请求参数的结构体
	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var schoolparam model.School

	json.Unmarshal([]byte(data), &schoolparam)

	err := h.schoolService.EditSchool(schoolparam, files)

	h.logger.Info("EditSchool", zap.Any("EditSchool", data))
	// 返回错误信息给客户端
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码
	resp.HandleSuccess(ctx, nil)
}

// 删除学校
func (h *schoolHandler) DelSchool(ctx *gin.Context) {
	var schoolparam model.School
	if err := ctx.ShouldBind(&schoolparam); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.schoolService.DelSchool(schoolparam)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}

// 查询学校及封面
func (h *schoolHandler) QuerySchoolAll(ctx *gin.Context) {

	data, err := h.schoolService.QuerySchoolAll()
	//记录获取的学校
	h.logger.Info("QuerySchoolAll ", zap.Any("QuerySchoolAll", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, data)
}
