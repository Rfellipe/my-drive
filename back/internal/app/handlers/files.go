package handlers

import (
	"fmt"
	"my-drive/internal/app/db"
	"my-drive/internal/pkg/lib"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleFileUpload(ctx *gin.Context) {
	database := db.DB.Connection
	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	dir := ctx.Query("dir")

	claims, error := lib.RetrieveJWTClaims(ctx)
	if error != nil {
		lib.RespondError(ctx, http.StatusForbidden, "No user ID found")
		return
	}

	if dir != "" {
		err := database.QueryRow(
			"SELECT id, owner, parent_id, name, type FROM nodes WHERE id=$1 AND owner=$1",
		).Scan(&dir)
		if err != nil {
			lib.RespondError(ctx, http.StatusInternalServerError, "Error connecting to database")
			return
		}
	}

	var fullPath string = ""
	for _, file := range files {
		if dir != "" {
			fullPath = fmt.Sprintf("%s/%s/%s", os.Getenv("DRIVE_ROOT"), claims.RootDirId, file.Filename)
		} else {
			fullPath = fmt.Sprintf("%s/%s/%s", os.Getenv("DRIVE_ROOT"), claims.RootDirId, file.Filename)
		}
		ctx.SaveUploadedFile(file, fullPath)
	}
	lib.Responder(ctx, http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)), nil)
}
