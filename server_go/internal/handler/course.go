package handler

import (
	v1 "ai_summary_project/api/v1"
	"ai_summary_project/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type CourseHandler interface {
	UploadCourse(ctx *gin.Context)
	GetCoursesByBookID(ctx *gin.Context)
	DeleteCourse(ctx *gin.Context)
	GetBooks(ctx *gin.Context)
}

type courseHandler struct {
	*Handler
	courseService service.CourseService // 服务层结构引入
}

func NewCourseHandler(
	handler *Handler,
	courseService service.CourseService,
) CourseHandler {
	return &courseHandler{
		Handler:       handler,
		courseService: courseService, // 服务层结构引入
	}
}

// POST /courses/upload
func (h *courseHandler) UploadCourse(ctx *gin.Context) {
	var req v1.UploadCourseRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.logger.Warn("UploadCourse", zap.Any("error", err))
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	h.logger.Info("Enter uploadCourse", zap.Any("BookID", req.BookID))
	var (
		mp4Header *multipart.FileHeader
		err       error
	)
	if req.Type == 1 {
		if mp4Header, err = ctx.FormFile("video"); err != nil {
			h.logger.Error("video get error", zap.Any("error", err))
			v1.HandleError(ctx, http.StatusBadRequest, err, nil)
			return
		}
	}
	mp3Header, err := ctx.FormFile("audio")
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	summary, err := h.courseService.GenerateSummaryFromAudio(ctx, mp3Header, req.Title)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	if req.Type == 1 {
		if err := h.courseService.CreateCourseWithFiles(
			ctx, req.BookID, req.Title, summary,
			mp4Header, mp3Header,
		); err != nil {
			h.logger.Error("UploadCourse", zap.Any("error", err))
			v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
			return
		}
	}

	v1.HandleSuccess(ctx, summary)
}

// GET /courses/:id
func (h *courseHandler) GetCoursesByBookID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	bookID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("GetCoursesByBookID Error!", zap.Any("error", err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrInvalidBookID, nil)
		return
	}
	h.logger.Info("GetCoursesByBookID", zap.Any("bookID", bookID))
	courses, err := h.courseService.GetCoursesByBookID(ctx, uint(bookID))
	if err != nil {
		h.logger.Error("GetCoursesByBookID Error!", zap.Any("error", err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrQueryCoursesFailed, nil)
		return
	}

	records := make([]v1.CourseRecord, 0, len(courses))
	for _, course := range courses {
		records = append(records, v1.CourseRecord{
			CourseID: course.ID,
			Title:    course.Title,
			Summary:  course.Summary,
			Video:    course.VideoPath,
			Record:   course.RecordPath,
			CreateAt: course.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	v1.HandleSuccess(ctx, v1.CourseListResponse{
		Total:   len(records),
		Records: records,
	})
}

// DELETE /courses/:id
func (h *courseHandler) DeleteCourse(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("DeleteCourse Error!", zap.Any("error", err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrInvalidCourseID, nil)
		return
	}
	h.logger.Info("DeleteCourse", zap.Any("id", id))
	if err := h.courseService.DeleteCourse(ctx, uint(id)); err != nil {
		h.logger.Error("DeleteCourse Error!", zap.Any("error", err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrDeleteCourseFailed, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// GET /books
func (h *courseHandler) GetBooks(ctx *gin.Context) {
	books, err := h.courseService.GetAllBooks(ctx)
	if err != nil {
		h.logger.Error("GetBooks Error!", zap.Any("error", err))
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrQueryBooksFailed, nil)
		return
	}
	h.logger.Info("Books Information")
	records := make([]v1.BookRecord, 0, len(books))
	for _, book := range books {
		records = append(records, v1.BookRecord{
			BookID:      book.ID,
			BookName:    book.BookName,
			CourseCount: book.CourseCount,
			CreateAt:    book.CreatedAt.Format("2006-1-2 15:04"),
			UpdateAt:    book.UpdatedAt.Format("2006-1-2 15:04"),
			DeleteAt:    formatNullableTime(book.DeletedAt),
		})
	}

	v1.HandleSuccess(ctx, v1.BookListResponse{
		Total:   len(records),
		Records: records,
	})
}

// 私有工具函数：格式化删除时间
func formatNullableTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format("2006-1-2 15:04")
}
