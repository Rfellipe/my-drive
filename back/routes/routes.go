package routes

import (
	_ "net/http"
	"github.com/gin-gonic/gin"
	"my-drive/handlers"
)

func Auth(router *gin.Engine) {
	router.POST("/auth/register", handlers.RegistrationHandler)
	router.POST("/auth/login", handlers.LoginHandler)
}

// func fileRoutes(router *gin.Engine) {
// 	router.GET("/file/upload", func(ctx *gin.Context) {
// 		ctx.JSON(http.StatusOK, gin.H{
// 			"req": "ok",
// 		})
// 	})
// }
