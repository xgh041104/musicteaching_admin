package server

import (
	"server_go/internal/handler"
	"server_go/internal/middleware"
	"server_go/pkg/log"

	"github.com/gin-gonic/gin"
)

func NewServerHTTP(
	logger *log.Logger,
	userHandler handler.UserHandler,
	adminHandler handler.AdminHandler,
	courseHandler handler.CourseHandler,
	resourcecategoryHandler handler.ResourceCategoryHandler,
	schoolHandler handler.SchoolHandler,
	courseCategoryHandler handler.CourseCategoryHandler,
	chapterHandler handler.ChapterHandler,
	sectionHandler handler.SectionHandler,
	commentHandler handler.CommentHandler,
	resourceHandler handler.ResourceHandler,
	giveScoreHandler handler.GiveScoreHandler,
	giveScoreThemeHandler handler.GiveScoreThemeHandler,
	giveScoreCategoryHandler handler.GiveScoreCategoryHandler,
	accountHandler handler.AccountHandler,
	studentUserHandler handler.StudentUserHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(
		middleware.CORSMiddleware(),
	)

	r.POST("/loginUser", userHandler.LoginUser)

	// user := r.Group("/user")  //普通用户接口
	// schooladmin := r.Group("/schooladmin") //学校管理员接口
	// administrator := r.Group("/administrator") //超级管理员接口

	//TODO
	// middleware.JwtAuthMiddleware(),  中间件,认证token合法性
	// 使用接口中间件  示例 : user.POST("/getArticleQuery",middleware.JwtAuthMiddleware(), userHandler.LoginUser)
	//  // 该函数调用middleware.JwtAuthMiddleware()  需要在我们函数的前面
	// user := r.Group("/manage", middleware.JwtAuthMiddleware()) //普通用户接口
	// {
	// }
	schooladmin := r.Group("/schooladmin", middleware.JwtAuthMiddleware()) //学校管理员接口
	{
		// 查询用户信息
		schooladmin.GET("/queryCommonUserBySchoolId", userHandler.GetUserBySchoolId)
		// 添加用户
		schooladmin.POST("/addCommonUser", userHandler.AddCommonUser)
		// 删除用户
		schooladmin.POST("/delCommonUser", userHandler.DelCommonUser)
		// 更新用户信息
		schooladmin.POST("/editCommonUser", userHandler.UpdateCommonUser)
		schooladmin.POST("/addCourse", courseHandler.AddCourse)
		schooladmin.POST("/editCourse", courseHandler.EditCourse)
		schooladmin.POST("/delCourse", courseHandler.DelCourse)

		//学生管理
		schooladmin.GET("/queryStudentByTeacherId", studentUserHandler.QueryStudentListByTeacherId)
		schooladmin.GET("/queryStudent", studentUserHandler.QueryStudentList)
		schooladmin.POST("/addStudent", studentUserHandler.AddStudent)
		schooladmin.POST("/addStudentBatch", studentUserHandler.AddStudentBatch)
		schooladmin.POST("/editStudent", studentUserHandler.EditStudent)
		schooladmin.POST("/delStudent", studentUserHandler.DelStudent)
	}

	administrator := r.Group("/administrator", middleware.JwtAuthMiddleware()) //超级管理员接口
	{
		// 查询学校管理员
		administrator.GET("/querySchoolAdmin", adminHandler.GetSchoolAdmin)
		// 添加学校管理员
		administrator.POST("/addSchoolAdmin", adminHandler.AddSchoolAdmin)
		// 删除学校管理员
		administrator.POST("/delSchoolAdmin", adminHandler.DelSchoolAdmin)
		// 更新学校管理员
		administrator.POST("/editSchoolAdmin", adminHandler.UpdateSchoolAdmin)

		administrator.POST("/addSchool", schoolHandler.AddSchool)
		administrator.POST("/editSchool", schoolHandler.EditSchool)
		administrator.POST("/delSchool", schoolHandler.DelSchool)
		administrator.GET("/querySchoolAll", schoolHandler.QuerySchoolAll)

		administrator.POST("/delComment", commentHandler.DelComment)
		administrator.GET("/queryComment", commentHandler.QueryComment)
		administrator.GET("/queryCommentByCommentId", commentHandler.QueryCommentByCommentId)

		administrator.POST("/addAccount", accountHandler.AddAccount)
		administrator.POST("/editAccount", accountHandler.EditAccount)
		administrator.POST("/delAccount", accountHandler.DelAccount)

		administrator.GET("/queryPaidCourse", courseHandler.QueryPaidCourse)
	}

	currency := r.Group("/currency", middleware.JwtAuthMiddleware()) //通用接口  三种用户都能访问
	{
		currency.POST("/editPwd", userHandler.EditUserPwd)

		currency.POST("/addCourse", courseHandler.AddCourse)
		currency.POST("/editCourse", courseHandler.EditCourse)
		currency.POST("/delCourse", courseHandler.DelCourse)
		currency.GET("/queryCourse", courseHandler.QueryCourse)
		currency.GET("/queryCourseByCategoryId", courseHandler.QueryCourseByCategoryId)
		currency.GET("/queryCourseDirectory", courseHandler.QueryCourseDirectory)
		currency.GET("/queryCourseCategoryTree", courseCategoryHandler.QueryCourseCategoryTree)
		currency.GET("/queryPublicCourse", courseHandler.QueryPublicCourse)
		currency.GET("/querySchoolCourse", courseHandler.QuerySchoolCourse)

		currency.POST("/addChapter", chapterHandler.AddChapter)

		currency.POST("/editChapter", chapterHandler.EditChapter)
		currency.POST("/delChapter", chapterHandler.DelChapter)
		currency.GET("/queryChapterByCourseId", chapterHandler.QueryChapterByCourseId)

		currency.POST("/editerUploadFile", sectionHandler.EditerUploadFile)
		currency.POST("/addSection", sectionHandler.AddSection)
		currency.POST("/delSection", sectionHandler.DelSection)
		currency.POST("/editSection", sectionHandler.EditSection)
		currency.GET("/querySectionBySectionId", sectionHandler.QuerySectionBySectionId)
		currency.GET("/querySectionByChapterId", sectionHandler.QuerySectionByChapterId)

		currency.POST("/addResource", resourceHandler.AddResource)
		currency.POST("/delResource", resourceHandler.DelResource)
		currency.GET("/queryResourceByResourceCategoryId", resourceHandler.QueryResourceByResourceCategoryId)
		currency.GET("/queryResourceById", resourceHandler.QueryResourceById)

		currency.GET("/queryResourceCategoryParentNodeByParentId", resourcecategoryHandler.QueryResourceCategoryParentNodeByParentId)
		currency.GET("/queryResourceCategoryChildNodesById", resourcecategoryHandler.QueryResourceCategoryChildNodesById)
		currency.GET("/queryResourceCategoryTree", resourcecategoryHandler.QueryResourceCategoryTree)

		currency.GET("/queryAllAccount", accountHandler.QueryAllAccount)
	}

	backstage := r.Group("/backstage", middleware.JwtAuthMiddleware()) //管理员和超管共用接口
	{

		backstage.POST("/addCourseCategory", courseCategoryHandler.AddCourseCategory)
		backstage.POST("/editCourseCategory", courseCategoryHandler.EditCourseCategory)
		backstage.POST("/delCourseCategory", courseCategoryHandler.DelCourseCategory)

		backstage.POST("/addResourceCategory", resourcecategoryHandler.AddResourceCategory)
		backstage.POST("/editResourceCategory", resourcecategoryHandler.EditResourceCategory)
		backstage.POST("/delResourceCategory", resourcecategoryHandler.DelResourceCategory)

	}

	commonuser := r.Group("/commonuser", middleware.JwtAuthMiddleware()) //普通用户
	{
		//添加评论
		commonuser.POST("/addComment", commentHandler.AddComment)
		commonuser.POST("/collectCourse", courseHandler.CollectCourse)
		commonuser.GET("/queryCourseIsCollected", courseHandler.QueryCourseIsCollected)
		commonuser.GET("/queryCollectCourse", courseHandler.QueryCollectCourse)

		commonuser.POST("/addGiveScoreCategory", giveScoreCategoryHandler.AddGiveScoreCategory)
		commonuser.POST("/editGiveScoreCategory", giveScoreCategoryHandler.EditGiveScoreCategory)
		commonuser.POST("/delGiveScoreCategory", giveScoreCategoryHandler.DelGiveScoreCategory)
		commonuser.GET("/queryGiveScoreCategory", giveScoreCategoryHandler.QueryGiveScoreCategory)
		commonuser.GET("/queryGiveScoreCategoryTree", giveScoreCategoryHandler.QueryGiveScoreCategoryTree)
		commonuser.GET("/queryGiveScoreCategoryChildNodesById", giveScoreCategoryHandler.QueryGiveScoreCategoryChildNodesById)
		commonuser.GET("/queryGiveScoreCategoryParentNodeByParentId", giveScoreCategoryHandler.QueryGiveScoreCategoryParentNodeByParentId)

		commonuser.POST("/addGiveScoreTheme", giveScoreThemeHandler.AddGiveScoreTheme)
		commonuser.POST("/editGiveScoreTheme", giveScoreThemeHandler.EditGiveScoreTheme)
		commonuser.POST("/delGiveScoreTheme", giveScoreThemeHandler.DelGiveScoreTheme)
		commonuser.GET("/queryGiveScoreThemeByGiveScoreCategoryId", giveScoreThemeHandler.QueryGiveScoreThemeByGiveScoreCategoryId)
		commonuser.GET("/queryGiveScoreTheme", giveScoreThemeHandler.QueryGiveScoreTheme)

		commonuser.POST("/editGiveScore", giveScoreHandler.EditGiveScore)
		commonuser.POST("/delGiveScore", giveScoreHandler.DelGiveSCore)
		commonuser.POST("/addGiveScore", giveScoreHandler.AddGiveScore)
		commonuser.GET("/queryGiveScore", giveScoreHandler.QueryGiveScore)
		commonuser.GET("/queryGiveTheacherScoreByGiveScoreThemeId", giveScoreHandler.QueryGiveTeacherScoreByGiveScoreThemeId)
		commonuser.GET("/queryGiveOneselfScoreByGiveScoreThemeId", giveScoreHandler.QueryGiveOneselfScoreByGiveScoreThemeId)
		commonuser.GET("/queryGiveMutualScoreByGiveScoreThemeId", giveScoreHandler.QueryGiveMutualScoreByGiveScoreThemeId)

		//查询个人课程
		commonuser.GET("/queryMyCourse", courseHandler.QueryMyCourse)
	}
	r.Static("/Resources/", "./Resources/")
	//TODO 发布需修改
	//r.Static("/Resources/", common.PublicPath+"/Resources/")
	return r
}
