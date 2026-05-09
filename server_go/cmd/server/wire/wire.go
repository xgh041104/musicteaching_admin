//go:build wireinject
// +build wireinject

package wire

import (
	"ai_summary_project/internal/handler"
	"ai_summary_project/internal/repository"
	"ai_summary_project/internal/server"
	"ai_summary_project/internal/service"
	"ai_summary_project/pkg/app"
	"ai_summary_project/pkg/jwt"
	"ai_summary_project/pkg/log"
	"ai_summary_project/pkg/server/http"
	"ai_summary_project/pkg/sid"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// 提取依赖集合
var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewCourseRepository,
	repository.NewBookRepository,
	repository.NewUsersRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewCourseService,
	// service.NewBookService,
	service.NewUsersService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewCourseHandler,
	// handler.NewBookHandler,
	handler.NewUsersHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
)

// 应用入口构建
func newApp(httpServer *http.Server) *app.App {
	return app.NewApp(
		app.WithServer(httpServer),
		app.WithName("course-server"),
	)
}

// Wire 入口
func NewWire(v *viper.Viper, logger *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
