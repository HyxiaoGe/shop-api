package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-api/goods-web/api/goods"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	zap.S().Info("初始化商品路由")
	{
		GoodsRouter.GET("", goods.List)
	}
}
