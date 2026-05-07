package handler

import (
	"encoding/json"
	"net/http"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type CourseCategoryHandler interface {
	AddCourseCategory(ctx *gin.Context)
	EditCourseCategory(ctx *gin.Context)
	DelCourseCategory(ctx *gin.Context)
	QueryCourseCategoryTree(ctx *gin.Context)
}

type courseCategoryHandler struct {
	*Handler
	courseCategoryService service.CourseCategoryService
}

func NewCourseCategoryHandler(handler *Handler, courseCategoryService service.CourseCategoryService) CourseCategoryHandler {
	return &courseCategoryHandler{
		Handler:               handler,
		courseCategoryService: courseCategoryService,
	}
}

func (h *courseCategoryHandler) AddCourseCategory(ctx *gin.Context) {

	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var courseCategoryparam model.CourseCategory

	json.Unmarshal([]byte(data), &courseCategoryparam)

	err := h.courseCategoryService.AddCourseCategory(courseCategoryparam, files)

	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return

	}
	resp.HandleSuccess(ctx, nil)

}

func (h *courseCategoryHandler) EditCourseCategory(ctx *gin.Context) {
	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var courseCategoryparam model.CourseCategory

	json.Unmarshal([]byte(data), &courseCategoryparam)

	err := h.courseCategoryService.EditCourseCategory(courseCategoryparam, files)

	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return

	}
	resp.HandleSuccess(ctx, nil)
}

func (h *courseCategoryHandler) DelCourseCategory(ctx *gin.Context) {
	var courseCategoryparam model.CourseCategory
	if err := ctx.ShouldBind(&courseCategoryparam); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.courseCategoryService.DelCourseCategory(courseCategoryparam)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}

func (h *courseCategoryHandler) QueryCourseCategoryTree(ctx *gin.Context) {
	var params struct {
		SchoolId int `form:"schoolId" ` // 请求参数 schoolId，
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.courseCategoryService.QueryCourseCategoryTree(params.SchoolId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, data)
}
