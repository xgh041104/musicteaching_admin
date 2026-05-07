package handler

import (
	"net/http"
	"server_go/internal/middleware"
	"server_go/internal/model"
	"server_go/internal/service"
	"server_go/pkg/helper/resp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler interface {
	LoginUser(ctx *gin.Context)
	EditUserPwd(cyx *gin.Context)
	GetUserBySchoolId(ctx *gin.Context)
	AddCommonUser(ctx *gin.Context)
	DelCommonUser(ctx *gin.Context)
	UpdateCommonUser(ctx *gin.Context)
}

type userHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(handler *Handler, userService service.UserService) UserHandler {
	return &userHandler{
		Handler:     handler,
		userService: userService,
	}
}

// 用户登录
func (h *userHandler) LoginUser(ctx *gin.Context) {
	var userparam model.LoginReq
	if err := ctx.ShouldBind(&userparam); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	user, err := h.userService.LoginUser(userparam)
	h.logger.Info("LoginUser", zap.Any("user", userparam))
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	if user == nil {
		resp.HandleError(ctx, http.StatusOK, 1, "账号或密码错误", nil)
		return
	}
	token, err := middleware.GetToken(user.UserAccount, user.UserType, user.SchoolId, user.UserId)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, "token获取错误", nil)
		return
	}

	resp.HandleLoginSuccess(ctx, user, token)
}

func (h *userHandler) EditUserPwd(ctx *gin.Context) {
	var user model.LoginUser
	if err := ctx.ShouldBind(&user); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.userService.EditUserPwd(user)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	resp.HandleSuccess(ctx, nil)
}

// 通过学校id进行用户查询
func (h *userHandler) GetUserBySchoolId(ctx *gin.Context) {
	//定义绑定请求参数的结构体
	var params struct {
		SchoolId int64 `form:"schoolId" binding:"required"` // 请求参数 schoolId，必须提供且不能为空
	}

	//将请求参数绑定到结构体
	if err := ctx.ShouldBind(&params); err != nil {
		//如果绑定参数失败，返回错误信息给客户端
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}
	//调用userService的GetUser方法获取用户信息
	commonuser, err := h.userService.GetUserBySchoolId(params.SchoolId)
	//记录获取的用户信息
	h.logger.Info("queryCommonUserBySchoolId", zap.Any("commonuser", commonuser))

	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码和用户信息
	resp.HandleSuccess(ctx, commonuser)
}

// 删除用户
func (h *userHandler) DelCommonUser(ctx *gin.Context) {
	//定义绑定请求参数的结构体
	var params struct {
		CommonUserId int64 `json:"userId" binding:"required"` //请求参数userId，必须提供且不能为空
	}

	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.userService.DelCommonUser(params.CommonUserId)
	//如果删除失败，返回错误信息给客户端
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	//返回成功状态码
	resp.HandleSuccess(ctx, nil)
}

// 编辑用户
func (h *userHandler) UpdateCommonUser(ctx *gin.Context) {

	var commonUser model.CommonUser
	if err := ctx.ShouldBind(&commonUser); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.userService.UpdateCommonUser(commonUser)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}

func (h *userHandler) AddCommonUser(ctx *gin.Context) {
	var commonUser model.CommonUser
	if err := ctx.ShouldBind(&commonUser); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	err := h.userService.AddCommonUser(commonUser)
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, nil)
}
