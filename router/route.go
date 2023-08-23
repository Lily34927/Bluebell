package router

import (
	"chapter4.1.bluebell/controller"
	"chapter4.1.bluebell/logger"
	"chapter4.1.bluebell/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式， cmd中不会有提示，否则是开发者模式，会有信息的提示
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")

	//注册业务路由
	v1.POST("/signup", controller.SighUpHandler)
	//登录
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middleware.JWTAuthMiddleware()) // 应用JWT认证中间件
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/posts", controller.GetPostListHandler)
		// 根据时间或分数获取帖子列表
		v1.GET("/posts2", controller.GetPostListHandler2)

		// 投票
		v1.POST("/vote", controller.PostVoteController)

	}
	//v1.GET("/ping", middleware.JWTAuthMiddleware(), func(c *gin.Context) { // 添加认证中间件
	//	// 如果是登录的用户，判断请求头中是否有 有效的JWT
	//	c.String(http.StatusOK, "pong")
	//})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
