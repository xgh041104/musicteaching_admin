package handler

import (
	v1 "ai_summary_project/api/v1"
	"ai_summary_project/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersHandler interface {
	GetUserList(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type usersHandler struct {
	*Handler
	usersService service.UsersService
}

func NewUsersHandler(
	handler *Handler,
	usersService service.UsersService,
) UsersHandler {
	return &usersHandler{
		Handler:      handler,
		usersService: usersService,
	}
}

func (h *usersHandler) GetUserList(ctx *gin.Context) {

}

func (h *usersHandler) Login(ctx *gin.Context) {
	req := new(v1.Login)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	user, err := h.usersService.Login(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, user)
}
