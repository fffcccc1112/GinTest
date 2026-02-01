package router

import (
	"github.com/gin-gonic/gin"
	"test/config"
	"test/internal/api/handler"
	"test/internal/api/middleware"
)

func NewRouter(cfg *config.Config,
	userHandler *handler.UserHandler) *gin.Engine {
	//初始化gin引擎
	gin.SetMode(cfg.ServerConfig.Mode)
	r := gin.New()
	//注册全局拦截器
	r.Use(middleware.Auth(cfg))
	r.Use(middleware.Logger())

	//健康检查
	r.GET("/healthy", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"code": 200,
			"msg":  "healthy",
		})
	})
	apiV1 := r.Group("/api")
	{
		apiV1.GET("/hello", userHandler.GetUserByID)

		apiV1.GET("/all", userHandler.GetALLUsers)
	}
	return r
}
