package server

import (
	"ai_summary_project/internal/handler"
	"ai_summary_project/internal/middleware"
	"ai_summary_project/pkg/log"
	"ai_summary_project/pkg/server/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	courseHandler handler.CourseHandler,
	UsersHandler handler.UsersHandler,
) *http.Server {
	if conf.GetString("env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	s.Use(middleware.CORSMiddleware())

	// 🎼 music teaching interface group
	music := s.Group("/musicTeaching")
	{
		music.POST("/course", courseHandler.UploadCourse)        // 添加课程
		music.DELETE("/course/:id", courseHandler.DeleteCourse)  // 删除课程
		music.GET("/book/:id", courseHandler.GetCoursesByBookID) // 某书下课程
		music.GET("/books", courseHandler.GetBooks)              // 获取所有书籍
	}

	s.Static("/static", "./static")
	s.POST("/login", UsersHandler.Login)
	return s
}
