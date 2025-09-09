package routes

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"my-drive/internal/app/handlers"
	"my-drive/internal/app/middleware"
)

func Auth(router *gin.Engine) {
	router.POST("/auth/register", handlers.RegistrationHandler)
	router.POST("/auth/login", handlers.LoginHandler)
}

func File(router *gin.Engine) {
	authorize := router.Group("/files")

	authorize.Use(middleware.JWTAuthorizeMiddleware())

	// Upload file
	authorize.POST("/upload", handlers.HandleFileUpload)

	// Download file
	authorize.GET("/:id", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"req": "ok",
		})
	})

	// Delete file
	authorize.DELETE("/:id", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"req": "ok",
		})
	})

	// List files
	authorize.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"req": "ok",
		})
	})
}
