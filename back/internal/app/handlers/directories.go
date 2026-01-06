package handlers

import (
	"fmt"
	"log"
	"my-drive/internal/app/db"
	"my-drive/internal/pkg/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DirForm struct {
	Name     string `json:"name" binding:"required"`
	ParentId string `json:"parent"`
}

type Directory struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func CreateDirectory(ctx *gin.Context) {
	var json DirForm

	claims, err := lib.RetrieveJWTClaims(ctx)
	if err != nil {
		log.Println("Error getting user claims", err)
		lib.RespondError(ctx, http.StatusForbidden, "Claims error")
		return
	}

	if err := ctx.ShouldBind(&json); err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	fmt.Println("ParentId:", json.ParentId)
	if err := lib.CreateDir(claims, json.Name, json.ParentId); err != nil {
		lib.RespondError(ctx, http.StatusInternalServerError, "ill need to handle it better")
		return
	}

	lib.Responder(ctx, http.StatusOK, "Create Directory Successfully", nil)
}

func FindDirectory(ctx *gin.Context) {
	database := db.DB.Connection
	var dir Directory
	var dirFiles []lib.FileInfo

	if err := ctx.ShouldBindUri(&dir); err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	claims, error := lib.RetrieveJWTClaims(ctx)
	if error != nil {
		lib.RespondError(ctx, http.StatusForbidden, "No user ID found")
		return
	}

	rows, err := database.Query(
		`
		SELECT
			name, type, size, created_at, updated_at
		FROM
			nodes
		WHERE
			parent_id = $1 AND owner = $2
		`,
		dir.ID,
		claims.Id,
	)
	if err != nil {
		fmt.Printf("%s", err)
		lib.RespondError(ctx, http.StatusInternalServerError, "Couldnt find directory")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var file lib.FileInfo
		if err := rows.Scan(&file.Name, &file.Type, &file.Size, &file.Created_at, &file.Updated_at); err != nil {
			fmt.Printf("%s", err)
			lib.RespondError(ctx, http.StatusInternalServerError, "Error reading files")
			return
		}

		dirFiles = append(dirFiles, file)
	}

	fmt.Println(dirFiles)
}
