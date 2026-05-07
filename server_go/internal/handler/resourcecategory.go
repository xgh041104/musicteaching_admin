package handler

import (
	"encoding/json"
	"net/http"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type ResourceCategoryHandler interface {
	AddResourceCategory(ctx *gin.Context)
	EditResourceCategory(ctx *gin.Context)
	DelResourceCategory(ctx *gin.Context)
	QueryResourceCategoryTree(ctx *gin.Context)

	QueryResourceCategoryParentNodeByParentId(ctx *gin.Context)
	QueryResourceCategoryChildNodesById(ctx *gin.Context)
}

type resourceCategoryHandler struct {
	*Handler
	resourceCategoryService service.ResourceCategoryService
}

func NewResourceCategoryHandler(handler *Handler, resourceCategoryService service.ResourceCategoryService) ResourceCategoryHandler {
	return &resourceCategoryHandler{
		Handler:                 handler,
		resourceCategoryService: resourceCategoryService,
	}
}

func (h *resourceCategoryHandler) AddResourceCategory(ctx *gin.Context) {
	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var resourcesCategoryparam model.ResourcesCategory

	json.Unmarshal([]byte(data), &resourcesCategoryparam)

	err := h.resourceCategoryService.AddResourceCategory(resourcesCategoryparam, files)

	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}

func (h *resourceCategoryHandler) EditResourceCategory(ctx *gin.Context) {
	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var resourcesCategoryparam model.ResourcesCategory

	json.Unmarshal([]byte(data), &resourcesCategoryparam)

	err := h.resourceCategoryService.EditResourceCategory(resourcesCategoryparam, files)

	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)

}
func (h *resourceCategoryHandler) DelResourceCategory(ctx *gin.Context) {

	var resourcesCategoryparam model.ResourcesCategory

	if err := ctx.ShouldBind(&resourcesCategoryparam); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	err := h.resourceCategoryService.DelResourceCategory(resourcesCategoryparam)

	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}
func (h *resourceCategoryHandler) QueryResourceCategoryTree(ctx *gin.Context) {
	var params struct {
		SchoolId int `form:"schoolId"` // 请求参数 学校id
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.resourceCategoryService.QueryResourceCategoryTree(params.SchoolId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), data)
		return
	}
	resp.HandleSuccess(ctx, data)
}

func (h *resourceCategoryHandler) QueryResourceCategoryParentNodeByParentId(ctx *gin.Context) {
	var params struct {
		ResourceCategoryParentId int `form:"resourceCategoryParentId"`
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.resourceCategoryService.QueryResourceCategoryParentNodeByParentId(params.ResourceCategoryParentId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), data)
		return
	}
	resp.HandleSuccess(ctx, data)
}
func (h *resourceCategoryHandler) QueryResourceCategoryChildNodesById(ctx *gin.Context) {
	var params struct {
		SchoolId           int `form:"schoolId"` // 请求参数 学校id
		ResourceCategoryId int `form:"resourceCategoryId"`
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.resourceCategoryService.QueryResourceCategoryChildNodesById(params.ResourceCategoryId, params.SchoolId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), data)
		return
	}
	resp.HandleSuccess(ctx, data)
}
