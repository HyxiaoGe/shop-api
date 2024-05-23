package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-api/user-web/api"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("初始化用户路由")
	{
		UserRouter.GET("list", api.GetUserList)
	}
}
