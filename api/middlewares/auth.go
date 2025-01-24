package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/metsasft/pkg"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}
		header := strings.Split(authHeader, " ")
		if strings.ToLower(header[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}
		claims, err := pkg.ValidateAccessToken(header[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
			ctx.Abort()
			return
		}
		ctx.Set("user_id", claims.UserID)
		ctx.Set("telegram_username", claims.TelegramUsername)
		ctx.Set("is_admin", claims.IsAdmin)
		ctx.Set("email", claims.Email)
		ctx.Set("expires_at", claims.ExpiresAt)
		ctx.Next()
	}
}
