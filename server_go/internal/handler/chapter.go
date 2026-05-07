package handler

import (
	"net/http"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/pkg/helper/resp"

	"github.com/gin-gonic/gin"
)

type ChapterHandler interface {
	AddChapter(ctx *gin.Context)
	EditChapter(ctx *gin.Context)
	DelChapter(ctx *gin.Context)
	QueryChapterByCourseId(ctx *gin.Context)
}

type chapterHandler struct {
	*Handler
	chapterService service.ChapterService
}

func NewChapterHandler(handler *Handler, chapterService service.ChapterService) ChapterHandler {
	return &chapterHandler{
		Handler:        handler,
		chapterService: chapterService,
	}
}

func (h *chapterHandler) AddChapter(ctx *gin.Context) {
	var chapter model.Chapter
	if err := ctx.ShouldBind(&chapter); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.chapterService.AddChapter(chapter)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}
func (h *chapterHandler) EditChapter(ctx *gin.Context) {
	var chapter model.Chapter
	if err := ctx.ShouldBind(&chapter); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.chapterService.EditChapter(chapter)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)

}
func (h *chapterHandler) DelChapter(ctx *gin.Context) {
	var chapter model.Chapter
	if err := ctx.ShouldBind(&chapter); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.chapterService.DelChapter(chapter)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)

}
func (h *chapterHandler) QueryChapterByCourseId(ctx *gin.Context) {
	var params struct {
		CourseId int `form:"courseId"` // 请求参数 courseId
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.chapterService.QueryChapterByCourseId(params.CourseId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), data)
		return
	}
	resp.HandleSuccess(ctx, data)

}
