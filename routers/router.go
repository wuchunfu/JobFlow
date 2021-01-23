package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/wuchunfu/JobFlow/controller/hostController"
	"github.com/wuchunfu/JobFlow/controller/loginLogController"
	"github.com/wuchunfu/JobFlow/controller/taskController"
	"github.com/wuchunfu/JobFlow/controller/taskLogController"
	"github.com/wuchunfu/JobFlow/controller/userController"
	"github.com/wuchunfu/JobFlow/controller/websocketController"
	"github.com/wuchunfu/JobFlow/middleware/cors"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Cors())
	// 设置session midddleware
	// 创建基于cookie的存储引擎，loginUser 参数是用于加密的密钥
	store := cookie.NewStore([]byte("loginUser"))
	// 设置 session 中间件，参数 mySession，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	router.Use(sessions.Sessions("mySession", store))
	gin.SetMode(gin.ReleaseMode)

	apiAuth := router.Group("/api/auth")
	{
		//apiAuth.POST("/register", userController.Register)
		//apiAuth.POST("/login", userController.Login)
		//apiAuth.GET("/info", auth.Auth(), userController.Info)
		//apiAuth.GET("/info", userController.Info)
		apiAuth.GET("/query", userController.Query)
		apiAuth.GET("/param/:param", userController.Param)
	}
	loginGroup := router.Group("/sys")
	{
		loginGroup.POST("/login", userController.Login)
		loginGroup.POST("/logout", userController.Logout)
	}
	router.GET("/sys/login/log/list", loginLogController.Index)
	taskGroup := router.Group("/sys/task")
	{
		taskGroup.GET("/list", taskController.Index)
		taskGroup.GET("/info/:taskId", taskController.Detail)
		taskGroup.GET("/hostAllList", taskController.HostAllList)
		taskGroup.POST("/save", taskController.Save)
		taskGroup.POST("/update", hostController.Update)
		taskGroup.POST("/delete", taskController.Delete)
		taskGroup.GET("/run/:taskId", taskController.Run)
		taskGroup.GET("/stop/:taskId", taskController.Stop)
		taskGroup.GET("/enable/:taskId", taskController.Enable)
		taskGroup.GET("/disable/:taskId", taskController.Disable)
		taskLogGroup := router.Group("/sys/task/log")
		{
			taskLogGroup.GET("/list", taskLogController.Index)
		}
	}
	hostGroup := router.Group("/sys/host")
	{
		hostGroup.GET("/list", hostController.Index)
		hostGroup.GET("/info/:hostId", hostController.Detail)
		hostGroup.POST("/save", hostController.Save)
		hostGroup.POST("/update", hostController.Update)
		hostGroup.POST("/delete", hostController.Delete)
		hostGroup.GET("/ping/:hostId", hostController.Ping)
	}
	userGroup := router.Group("/sys/user")
	{
		userGroup.GET("/list", userController.Index)
		userGroup.GET("/info/:userId", userController.Detail)
		userGroup.POST("/save", userController.Save)
		userGroup.POST("/update", userController.Update)
		userGroup.POST("/changePassword", userController.ChangePassword)
		userGroup.POST("/changeLoginPassword", userController.ChangeLoginPassword)
		userGroup.POST("/delete", userController.Delete)
	}

	router.GET("/websocket", websocketController.LogTail)

	//userGroup := router.Group("/api/auth/user")
	//{
	//	userGroup.GET("", userController.Index)
	//	userGroup.GET("/:id", userController.Detail)
	//	userGroup.POST("/save", userController.Save)
	//	userGroup.POST("/update", userController.Update)
	//	userGroup.POST("/updatePassword", userController.UpdatePassword)
	//	userGroup.POST("/updateLoginPassword", userController.UpdateLoginPassword)
	//	userGroup.DELETE("/delete", userController.Delete)
	//}
	router.NoRoute(userController.NotFound)
	return router
}
