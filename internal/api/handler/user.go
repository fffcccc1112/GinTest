package handler

import (
	"github.com/gin-gonic/gin"
	"test/internal/service"
	"test/pkg/logger"
)

// 接口处理器，相当于控制层
// user控制层
type UserHandler struct {
	UserService *service.UserService
}

// 创建
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}
func (h *UserHandler) GetUserByID(c *gin.Context) {
	user, err := h.UserService.GetUserByID(c.Request.Context(), 122)
	if err != nil {
		logger.Error("查询失败!")
	}
	c.JSON(200,
		gin.H{
			"code": 200,
			"data": user,
		})
}
func (h *UserHandler) GetALLUsers(c *gin.Context) {
	users, err := h.UserService.GetALLUsers()
	if err != nil {
		c.JSON(500,
			gin.H{
				"code": 500,
				"data": []string{},
			})
	}
	c.JSON(200,
		gin.H{
			"code": 200,
			"data": users,
		})
}
