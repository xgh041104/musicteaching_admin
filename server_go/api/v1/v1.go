package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	resp := Response{Code: errorCodeMap[ErrSuccess], Message: ErrSuccess.Error(), Data: data}
	if _, ok := errorCodeMap[ErrSuccess]; !ok {
		resp = Response{Code: http.StatusOK, Message: "", Data: data}
	}
	ctx.JSON(http.StatusOK, resp)
}

func HandleError(ctx *gin.Context, httpCode int, err error, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	var (
		code int
		msg  string = err.Error()
	)
	switch e := err.(type) {
	case validator.ValidationErrors:
		code = errorCodeMap[ErrBadRequest]
		// ✅ 构建字段错误提示
		var buf strings.Builder
		for _, fieldErr := range e {
			buf.WriteString(fmt.Sprintf("[%s] 不合法; ", fieldErr.Field()))
		}
		msg = buf.String()
	default:
		if val, ok := errorCodeMap[err]; ok {
			code = val
		} else {
			code = 500
		}
	}

	resp := Response{Code: code, Message: msg, Data: data}
	ctx.JSON(httpCode, resp)
}

type Error struct {
	Code    int
	Message string
}

var errorCodeMap = map[error]int{}

func newError(code int, msg string) error {
	err := errors.New(msg)
	errorCodeMap[err] = code
	return err
}
func (e Error) Error() string {
	return e.Message
}
