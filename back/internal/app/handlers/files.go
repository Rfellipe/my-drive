package handlers

import (
	"fmt"
	"my-drive/internal/pkg/lib"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleFileUpload(ctx *gin.Context) {
	// Multipart form
	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	for _, file := range files {
		fmt.Println(file.Filename)

		fullPath := fmt.Sprintf("%s/%s", os.Getenv("DRIVE_ROOT"), file.Filename)
		fmt.Println(fullPath)
		// Upload the file to specific dst.
		ctx.SaveUploadedFile(file, fullPath)
	}
	lib.Responder(ctx, http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)), nil)
}
