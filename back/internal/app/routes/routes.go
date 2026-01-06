package routes

import (
	"github.com/gin-gonic/gin"
	"my-drive/internal/app/handlers"
	"my-drive/internal/app/middleware"
	"net/http"
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

func Dir(router *gin.Engine) {
	authorize := router.Group("/dir")

	authorize.Use(middleware.JWTAuthorizeMiddleware())

	// Get folder
	authorize.GET("/:id", handlers.FindDirectory)

	// Update folder
	authorize.PUT("/:id", func(ctx *gin.Context) {})

	// Create folder
	authorize.POST("/", handlers.CreateDirectory)

	// Delete folder
	authorize.DELETE("/:id", func(ctx *gin.Context) {})
}
