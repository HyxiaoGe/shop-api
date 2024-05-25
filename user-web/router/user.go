package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-api/user-web/api"
	"shop-api/user-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("初始化用户路由")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.LoginByPassword)
		UserRouter.POST("register", api.Register)
	}
}
