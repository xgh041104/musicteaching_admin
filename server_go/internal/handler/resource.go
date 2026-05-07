package handler

import (
	"encoding/json"
	"net/http"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type ResourceHandler interface {
	AddResource(ctx *gin.Context)
	DelResource(ctx *gin.Context)
	QueryResourceByResourceCategoryId(ctx *gin.Context)
	QueryResourceById(ctx *gin.Context)
}

type resourceHandler struct {
	*Handler
	resourceService service.ResourceService
}

func NewResourceHandler(handler *Handler, resourceService service.ResourceService) ResourceHandler {
	return &resourceHandler{
		Handler:         handler,
		resourceService: resourceService,
	}
}

func (h *resourceHandler) AddResource(ctx *gin.Context) {
	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var resourcesparam model.Resources

	json.Unmarshal([]byte(data), &resourcesparam)

	err := h.resourceService.AddResource(resourcesparam, files)

	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}

func (h *resourceHandler) DelResource(ctx *gin.Context) {
	var resource model.Resources
	if err := ctx.ShouldBind(&resource); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	err := h.resourceService.DelResource(resource)

	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}

func (h *resourceHandler) QueryResourceByResourceCategoryId(ctx *gin.Context) {

	var params struct {
		SchoolId             int `form:"schoolId"`
		LecturerCommonUserId int `form:"lecturerCommonUserId"`
		ResourceCategoryId   int `form:"resourceCategoryId"`
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	data, err := h.resourceService.QueryResourceByResourceCategoryId(params.SchoolId, params.LecturerCommonUserId, params.ResourceCategoryId)

	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, data)
}

func (h *resourceHandler) QueryResourceById(ctx *gin.Context) {

	var params struct {
		ResourceId int `form:"resourceId"`
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	data, err := h.resourceService.QueryResourceById(params.ResourceId)

	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, data)
}
