package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-api/goods-web/middlewares"
	"shop-api/goods-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
			"msg":     "ok",
		})
	})
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/goods/v1")
	router.InitGoodsRouter(ApiGroup)
	return Router
}
