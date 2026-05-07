package handler

import (
	"encoding/json"
	"net/http"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/pkg/helper/resp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AccountHandler interface {
	AddAccount(ctx *gin.Context)
	EditAccount(ctx *gin.Context)
	DelAccount(ctx *gin.Context)
	QueryAllAccount(ctx *gin.Context)
}

type accountHandler struct {
	*Handler
	accountService service.AccountService
}

func NewAccountHandler(handler *Handler, accountService service.AccountService) AccountHandler {
	return &accountHandler{
		Handler:        handler,
		accountService: accountService,
	}
}

func (h *accountHandler) AddAccount(ctx *gin.Context) {

	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var accountparam model.Account

	json.Unmarshal([]byte(data), &accountparam)

	err := h.accountService.AddAccount(accountparam, files)
	//记录获取的公众号
	h.logger.Info("AddAccount", zap.Any("AddAccount", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, nil)
}

// 修改公众号
func (h *accountHandler) EditAccount(ctx *gin.Context) {
	//定义绑定请求参数的结构体
	files, _ := ctx.MultipartForm() // 获取fromdata的文件

	data, _ := ctx.GetPostForm("data") // 获取fromdata的文件

	var accountparam model.Account

	json.Unmarshal([]byte(data), &accountparam)

	err := h.accountService.EditAccount(accountparam, files)

	h.logger.Info("EditAccount", zap.Any("EditAccount", data))
	// 返回错误信息给客户端
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码
	resp.HandleSuccess(ctx, nil)
}

// 删除公众号
func (h *accountHandler) DelAccount(ctx *gin.Context) {
	var accountparam model.Account
	if err := ctx.ShouldBind(&accountparam); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.accountService.DelAccount(accountparam)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}

// 查询公众号及封面
func (h *accountHandler) QueryAllAccount(ctx *gin.Context) {

	data, err := h.accountService.QueryAllAccount()
	//记录获取的公众号
	h.logger.Info("QueryAllAccount ", zap.Any("QueryAllAccount", data))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, data)
}
