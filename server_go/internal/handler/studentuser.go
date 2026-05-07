package handler

import (
	"encoding/json"
	"net/http"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type StudentUserHandler interface {
	QueryStudentListByTeacherId(ctx *gin.Context)
	QueryStudentList(ctx *gin.Context)
	AddStudent(ctx *gin.Context)
	AddStudentBatch(ctx *gin.Context)
	EditStudent(ctx *gin.Context)
	DelStudent(ctx *gin.Context)
}

type studentUserHandler struct {
	*Handler
	studentUserService service.StudentUserService
}

func NewStudentUserHandler(handler *Handler, studentUserService service.StudentUserService) StudentUserHandler {
	return &studentUserHandler{
		Handler:            handler,
		studentUserService: studentUserService,
	}
}

func (h *studentUserHandler) QueryStudentList(ctx *gin.Context) {
	value, exists := ctx.Get("schoolId")
	if !exists {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	schoolId, ok := value.(int)
	if !ok {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	students, err := h.studentUserService.QueryStudentList(int64(schoolId))
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, students)
}

// 通过老师id进行用户查询
func (h *studentUserHandler) QueryStudentListByTeacherId(ctx *gin.Context) {
	var req struct {
		TeacherId int64 `form:"teacherId" binding:"required"`
	}

	if err := ctx.ShouldBind(&req); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	value, exists := ctx.Get("schoolId")
	if !exists {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	schoolId, ok := value.(int)
	if !ok {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	students, err := h.studentUserService.QueryStudentListByTeacherId(int64(req.TeacherId), int64(schoolId))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, students)
}

// 删除用户
func (h *studentUserHandler) DelStudent(ctx *gin.Context) {
	var req struct {
		StudentId int64 `json:"studentId" binding:"required"`
	}

	if err := ctx.ShouldBind(&req); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	value, exists := ctx.Get("schoolId")
	if !exists {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	schoolId, ok := value.(int)
	if !ok {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	err := h.studentUserService.DelStudent(req.StudentId, int64(schoolId))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}

// 编辑用户
func (h *studentUserHandler) EditStudent(ctx *gin.Context) {
	var req struct {
		StudentId     int    `json:"studentId" binding:"required"`
		StudentName   string `json:"studentName" binding:"required"`
		ClassName     string `json:"className" binding:"required"`
		TeacherIdList []int  `json:"teacherIdList" binding:"required"`
	}

	if err := ctx.ShouldBind(&req); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	var student model.Student
	student.StudentId = req.StudentId
	student.StudentName = req.StudentName
	student.ClassName = req.ClassName
	student.TeacherIdList = req.TeacherIdList

	err := h.studentUserService.EditStudent(&student)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}

func (h *studentUserHandler) AddStudent(ctx *gin.Context) {
	var req struct {
		StudentAccount string `json:"studentAccount" binding:"required"`
		StudentPwd     string `json:"studentPwd" binding:"required"`
		StudentName    string `json:"studentName" binding:"required"`
		ClassName      string `json:"className" binding:"required"`
		TeacherId      int64  `json:"teacherId" binding:"required"`
	}

	if err := ctx.ShouldBind(&req); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	value, exists := ctx.Get("schoolId")
	if !exists {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	schoolId, ok := value.(int)
	if !ok {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	student := &model.StudentUser{
		StudentAccount: req.StudentAccount,
		StudentPwd:     req.StudentPwd,
		StudentName:    req.StudentName,
		ClassName:      req.ClassName,
		SchoolId:       schoolId,
	}

	err := h.studentUserService.AddStudent(student, req.TeacherId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}

func (h *studentUserHandler) AddStudentBatch(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("files")
	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	var req struct {
		TeacherId int64 `json:"teacherId" binding:"required"`
	}

	data, _ := ctx.GetPostForm("data")
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	value, exists := ctx.Get("schoolId")
	if !exists {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	schoolId, ok := value.(int)
	if !ok {
		resp.HandleError(ctx, http.StatusBadRequest, 1, "参数错误", nil)
		return
	}

	err = h.studentUserService.AddStudentBatch(fileHeader, req.TeacherId, int64(schoolId))
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}
