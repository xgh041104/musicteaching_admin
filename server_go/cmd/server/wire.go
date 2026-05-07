//go:build wireinject
// +build wireinject

package main

import (
	"server_go/internal/handler"
	"server_go/internal/repository"
	"server_go/internal/server"
	"server_go/internal/service"
	"server_go/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var ServerSet = wire.NewSet(server.NewServerHTTP)

var RepositorySet = wire.NewSet(
	repository.NewBaseDb,
	repository.NewBaseRepository,
	repository.NewUserRepository,
	repository.NewAdminRepository,
	repository.NewCourseRepository,
	repository.NewResourceCategoryRepository,
	repository.NewSchoolRepository,
	repository.NewCourseCategoryRepository,
	repository.NewChapterRepository,
	repository.NewSectionRepository,
	repository.NewCommentRepository,
	repository.NewResourceRepository,
	repository.NewGiveScoreRepository,
	repository.NewGiveScoreThemeRepository,
	repository.NewGiveScoreCategoryRepository,
	repository.NewAccountRepository,
	repository.NewStudentUserRepository,
)
var ServiceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewAdminService,
	service.NewCourseService,
	service.NewResourceCategoryService,
	service.NewSchoolService,
	service.NewCourseCategoryService,
	service.NewChapterService,
	service.NewSectionService,
	service.NewCommentService,
	service.NewResourceService,
	service.NewGiveScoreService,
	service.NewGiveScoreThemeService,
	service.NewGiveScoreCategoryService,
	service.NewAccountService,
	service.NewStudentUserService,
)

var HandlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewAdminHandler,
	handler.NewCourseHandler,
	handler.NewResourceCategoryHandler,
	handler.NewSchoolHandler,
	handler.NewCourseCategoryHandler,
	handler.NewChapterHandler,
	handler.NewSectionHandler,
	handler.NewCommentHandler,
	handler.NewResourceHandler,
	handler.NewgiveScoreHandler,
	handler.NewgiveScoreThemeHandler,
	handler.NewgiveScoreCategoryHandler,
	handler.NewAccountHandler,
	handler.NewStudentUserHandler,
)

func newApp(*viper.Viper, *log.Logger) (*gin.Engine, func(), error) {
	panic(wire.Build(
		ServerSet,
		RepositorySet,
		ServiceSet,
		HandlerSet,
	))
}
