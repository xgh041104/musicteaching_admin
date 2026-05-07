package handler

import (
	"net/http"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/pkg/helper/resp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminHandler interface {
	GetSchoolAdmin(ctx *gin.Context)
	AddSchoolAdmin(ctx *gin.Context)
	DelSchoolAdmin(ctx *gin.Context)
	UpdateSchoolAdmin(ctx *gin.Context)
}

type adminHandler struct {
	*Handler
	adminService service.AdminService
}

func NewAdminHandler(handler *Handler, adminService service.AdminService) AdminHandler {
	return &adminHandler{
		Handler:      handler,
		adminService: adminService,
	}
}

// 添加学校管理员
func (h *adminHandler) AddSchoolAdmin(ctx *gin.Context) {
	var admin model.Admin
	if err := ctx.ShouldBind(&admin); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.adminService.AddSchoolAdmin(admin)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}

// 获取学校管理员信息
func (h *adminHandler) GetSchoolAdmin(ctx *gin.Context) {

	//查询学校管理员
	schooladmin, err := h.adminService.GetSchoolAdmin()

	//记录获取的管理员信息
	h.logger.Info("querySchoolAdmin", zap.Any("schooladmin", schooladmin))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, schooladmin)
}

// 删除学校管理员
func (h *adminHandler) DelSchoolAdmin(ctx *gin.Context) {
	//定义绑定请求参数的结构体
	var params struct {
		AdminId int64 `json:"userId" binding:"required"` //请求参数userId，必须提供且不能为空
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.adminService.DelSchoolAdmin(params.AdminId)
	//如果删除失败，返回错误信息给客户端
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码
	resp.HandleSuccess(ctx, nil)
}

// 修改学校管理员信息
func (h *adminHandler) UpdateSchoolAdmin(ctx *gin.Context) {
	var admin model.Admin
	if err := ctx.ShouldBind(&admin); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.adminService.UpdateSchoolAdmin(admin)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}
