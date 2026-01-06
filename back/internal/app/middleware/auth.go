package middleware

import (
	_ "fmt"
	"my-drive/internal/pkg/lib"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthorizeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			lib.RespondError(ctx, http.StatusUnauthorized, "Missing Authorization header")
			ctx.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			lib.RespondError(ctx, http.StatusUnauthorized, "Invalid Authorization header format")
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		claims, err := lib.ValidateJWT(tokenString)
		if err != nil {
			lib.RespondError(ctx, http.StatusUnauthorized, "Invalid token")
			ctx.Abort()
			return
		}

		exp, err := claims.Claims.GetExpirationTime()
		if err != nil || exp == nil || exp.Time.Before(time.Now()) {
			lib.RespondError(ctx, http.StatusUnauthorized, "Token expired")
			ctx.Abort()
			return
		}

		sub, ok := claims.Claims.(jwt.MapClaims)["sub"].(map[string]any)
		if !ok {
			lib.RespondError(ctx, http.StatusUnauthorized, "Invalid subject claim")
			ctx.Abort()
			return
		}

		userID, _ := sub["id"].(string)
		email, _ := sub["email"].(string)
		status, _ := sub["status"].(string)
		rootDir, _ := sub["rootDirId"].(string)
		loginAttempts, _ := sub["loginAttempts"].(float64)

		ctx.Set("userId", userID)
		ctx.Set("userEmail", email)
		ctx.Set("userStatus", status)
		ctx.Set("userRootDir", rootDir)
		ctx.Set("userLoginAttempts", int(loginAttempts))
		ctx.Set("claims", claims.Claims)

		ctx.Next()
	}
}
