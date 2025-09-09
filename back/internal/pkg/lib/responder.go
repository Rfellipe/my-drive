package lib

import "github.com/gin-gonic/gin"

type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Responder(ctx *gin.Context, code int, message string, data any) {
	ctx.JSON(code, ApiResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func RespondError(ctx *gin.Context, code int, message string) {
	Responder(ctx, code, message, nil)
}

