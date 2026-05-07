package handler

import (
	"encoding/json"
	"net/http"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type CourseHandler interface {
	AddCourse(ctx *gin.Context)
	EditCourse(ctx *gin.Context)
	DelCourse(ctx *gin.Context)
	QueryCourse(ctx *gin.Context)
	QueryCourseByCategoryId(ctx *gin.Context)
	QueryCourseDirectory(ctx *gin.Context)
	CollectCourse(ctx *gin.Context)
	QueryCourseIsCollected(ctx *gin.Context)
	QueryPaidCourse(ctx *gin.Context)
	QueryCollectCourse(ctx *gin.Context)
	QueryPublicCourse(ctx *gin.Context)
	QuerySchoolCourse(ctx *gin.Context)
	QueryMyCourse(ctx *gin.Context)
}

type courseHandler struct {
	*Handler
	courseService service.CourseService
}

func NewCourseHandler(handler *Handler, courseService service.CourseService) CourseHandler {
	return &courseHandler{
		Handler:       handler,
		courseService: courseService,
	}
}

func (h *courseHandler) AddCourse(ctx *gin.Context) {

	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var courseparam model.Course

	json.Unmarshal([]byte(data), &courseparam)

	err := h.courseService.AddCourse(courseparam, files)

	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return

	}
	resp.HandleSuccess(ctx, nil)

}

func (h *courseHandler) EditCourse(ctx *gin.Context) {
	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var courseparam model.Course

	json.Unmarshal([]byte(data), &courseparam)

	err := h.courseService.EditCourse(courseparam, files)

	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return

	}
	resp.HandleSuccess(ctx, nil)
}

func (h *courseHandler) DelCourse(ctx *gin.Context) {
	var courseparam model.Course
	if err := ctx.ShouldBind(&courseparam); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.courseService.DelCourse(courseparam)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}

func (h *courseHandler) QueryCourse(ctx *gin.Context) {
	var params struct {
		SchoolId             int `form:"schoolId" ` // 请求参数 schoolId，
		LecturerCommonUserId int `form:"lecturerCommonUserId" `
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.courseService.QueryCourse(params.SchoolId, params.LecturerCommonUserId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), data)
		return
	}
	resp.HandleSuccess(ctx, data)
}

func (h *courseHandler) QueryCourseDirectory(ctx *gin.Context) {
	var params struct {
		CourseId   int `form:"courseId" ` // 请求参数 schoolId，
		CommUserId int `form:"commUserId" `
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	schoolIdInterface, _ := ctx.Get("schoolId")
	schoolId, _ := schoolIdInterface.(int)

	data, err := h.courseService.QueryCourseDirectory(params.CourseId, params.CommUserId, schoolId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), data)
		return
	}
	resp.HandleSuccess(ctx, data)
}

func (h *courseHandler) QueryCourseByCategoryId(ctx *gin.Context) {
	var params struct {
		SchoolId             int    `form:"schoolId" ` // 请求参数 schoolId
		TrueSchoolId         int    `form:"trueSchoolId" `
		LecturerCommonUserId int    `form:"lecturerCommonUserId" `
		Page                 int    `form:"page" `
		PageSize             int    `form:"pageSize" `
		TitleWords           string `form:"titleWords" `
		CategoryId           int    `form:"categoryId" `
		UserId               int    `form:"userId" `
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.courseService.QueryCourseByCategoryId(params.SchoolId, params.LecturerCommonUserId, params.Page, params.PageSize, params.CategoryId, params.TitleWords, params.TrueSchoolId, params.UserId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), data)
		return
	}
	resp.HandleSuccess(ctx, data)
}

func (h *courseHandler) CollectCourse(ctx *gin.Context) {
	var params struct {
		UserId    int  `json:"userId"`
		CourseId  int  `json:"courseId"`
		IsCollect bool `json:"isCollect"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.courseService.CollectCourse(params.CourseId, params.UserId, params.IsCollect)

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}

func (h *courseHandler) QueryCourseIsCollected(ctx *gin.Context) {
	var params struct {
		UserId   int `form:"userId"`
		CourseId int `form:"courseId"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	isCollected, err := h.courseService.QueryCourseIsCollected(params.CourseId, params.UserId)

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, isCollected)
}

func (h *courseHandler) QueryPaidCourse(ctx *gin.Context) {
	data, err := h.courseService.QueryPaidCourse()
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, data)
}

func (h *courseHandler) QueryCollectCourse(ctx *gin.Context) {
	var params struct {
		UserId     int    `form:"userId"`
		SchoolId   int    `form:"schoolId"`
		Page       int    `form:"page"`
		PageSize   int    `form:"pageSize"`
		TitleWords string `form:"titleWords"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.courseService.QueryCollectCourse(params.UserId, params.SchoolId, params.Page, params.PageSize, params.TitleWords)

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, data)
}

// 查询个人课程(含收藏课程)
func (h *courseHandler) QueryMyCourse(ctx *gin.Context) {
	var params struct {
		UserId     int    `form:"userId"`
		SchoolId   int    `form:"schoolId"`
		Page       int    `form:"page"`
		PageSize   int    `form:"pageSize"`
		TitleWords string `form:"titleWords"`
		CategoryId int    `form:"categoryId"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.courseService.QueryMyCourse(params.UserId, params.SchoolId, params.Page, params.PageSize, params.TitleWords, params.CategoryId)

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, data)
}

// 查询公共课程(不含收费)
func (h *courseHandler) QueryPublicCourse(ctx *gin.Context) {
	var params struct {
		UserId     int    `form:"userId"`
		Page       int    `form:"page"`
		PageSize   int    `form:"pageSize"`
		TitleWords string `form:"titleWords"`
		CategoryId int    `form:"categoryId"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.courseService.QueryPublicCourse(params.UserId, params.Page, params.PageSize, params.TitleWords, params.CategoryId)

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, data)
}

// 查询学校课程(含付费公共课程)
func (h *courseHandler) QuerySchoolCourse(ctx *gin.Context) {
	var params struct {
		UserId     int    `form:"userId"`
		SchoolId   int    `form:"schoolId"`
		Page       int    `form:"page"`
		PageSize   int    `form:"pageSize"`
		TitleWords string `form:"titleWords"`
		CategoryId int    `form:"categoryId"`
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.courseService.QuerySchoolCourse(params.UserId, params.SchoolId, params.Page, params.PageSize, params.TitleWords, params.CategoryId)

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, data)
}
