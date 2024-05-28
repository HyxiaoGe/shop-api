package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-api/goods-web/models"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentClaims := claims.(*models.CustomClaims)
		if currentClaims.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "当前无权限访问",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
