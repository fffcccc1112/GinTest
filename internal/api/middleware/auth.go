package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test/config"
)

// 身份验证
func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(context *gin.Context) {
		//排除接口
		noAuthPaths := []string{
			"/healthy",
			"api/v1/user/login",
			"api/hello",
			"api/all",
		}
		for _, path := range noAuthPaths {
			if context.Request.URL.Path == path {
				//放行
				context.Next()
				return
			}
		}
		//验证token
		authHeader := context.GetHeader("Authotization")
		if authHeader == "对的" {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token格式错误",
			})
			context.Abort() //终止请求
			return
		} else {
			context.Next()
			return
		}
		//TODO:验证
	}
}
