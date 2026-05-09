package middleware

import (
	"ai_summary_project/pkg/log"
	"bytes"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"strings"
	"time"
)

func RequestLogMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 忽略大文件接口，防止日志爆炸
		if strings.Contains(ctx.Request.URL.Path, "/upload") ||
			strings.Contains(ctx.Request.URL.Path, "/course") ||
			strings.HasPrefix(ctx.GetHeader("Content-Type"), "multipart/form-data") {
			ctx.Next()
			return
		}
		uuid, err := random.UUIdV4()
		if err != nil {
			return
		}
		trace := cryptor.Md5String(uuid)
		logger.WithValue(ctx, zap.String("trace", trace))
		logger.WithValue(ctx, zap.String("request_method", ctx.Request.Method))
		logger.WithValue(ctx, zap.Any("request_headers", ctx.Request.Header))
		logger.WithValue(ctx, zap.String("request_url", ctx.Request.URL.String()))

		// 仅打印最多 2 KB 请求体
		if ctx.Request.Body != nil {
			bodyBytes, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			logger.WithValue(ctx, zap.String("request_params", truncate(string(bodyBytes), 2048)))
		}

		logger.WithContext(ctx).Info("Request")
		ctx.Next()
	}
}

func ResponseLogMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		start := time.Now()
		ctx.Next()
		duration := time.Since(start)

		respStr := blw.body.String()
		logger.WithContext(ctx).Info("Response",
			zap.String("response_body", truncate(respStr, 4096)),
			zap.String("duration", duration.String()))
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "...(truncated)"
}
