package handler

import (
	"net/http"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/pkg/helper/resp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommentHandler interface {
	AddComment(ctx *gin.Context)
	DelComment(ctx *gin.Context)
	QueryCommentByCommentId(ctx *gin.Context)
	QueryComment(ctx *gin.Context)
}

type commentHandler struct {
	*Handler
	commentService service.CommentService
}

func NewCommentHandler(handler *Handler, commentService service.CommentService) CommentHandler {
	return &commentHandler{
		Handler:        handler,
		commentService: commentService,
	}
}

func (h *commentHandler) AddComment(ctx *gin.Context) {
	var comment model.Comment
	if err := ctx.ShouldBind(&comment); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.commentService.AddComment(comment)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}

func (h *commentHandler) DelComment(ctx *gin.Context) {
	//定义绑定请求参数的结构体
	var params struct {
		CommentId int64 `json:"commentId" binding:"required"` //请求参数commentId，必须提供且不能为空
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.commentService.DelComment(params.CommentId)
	//如果删除失败，返回错误信息给客户端
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码
	resp.HandleSuccess(ctx, nil)
}

func (h *commentHandler) QueryComment(ctx *gin.Context) {

	data, err := h.commentService.QueryComment()
	//记录获取的评论
	h.logger.Info("QueryComment ", zap.Any("QueryComment", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, data)
}

func (h *commentHandler) QueryCommentByCommentId(ctx *gin.Context) {

	var params struct {
		CommentId int64 `form:"commentId" binding:"required"` //请求参数commentId，必须提供且不能为空
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	data, err := h.commentService.QueryCommentByCommentId(params.CommentId)
	//记录获取的评论
	h.logger.Info("QueryComment ", zap.Any("QueryComment", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, data)
}
