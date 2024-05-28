package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-api/goods-web/api/goods"
	"shop-api/goods-web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	zap.S().Info("初始化商品路由")
	{
		GoodsRouter.GET("", goods.List)
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)
		GoodsRouter.GET("/:id", goods.Detail)
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)
		GoodsRouter.GET("/:id/stocks", goods.Stocks)
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)

	}
}
