package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, ok := ctx.Get("is_admin")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}
		isAdmin := value.(bool)
		if !isAdmin {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "admin access is required"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
