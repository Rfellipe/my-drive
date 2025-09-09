package handlers

import (
	"fmt"
	"my-drive/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleFileUpload(ctx *gin.Context) {
	// Multipart form
	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	userId, _ := ctx.Get("userId")

	fmt.Printf("User %v is uploading %d files\n", userId, len(files))

	// for _, file := range files {
	// 	fmt.Println(file.Filename)

	// 	// Upload the file to specific dst.
	// 	// ctx.SaveUploadedFile(file, "./files/"+file.Filename)
	// }
	lib.Responder(ctx, http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)), nil)
}
