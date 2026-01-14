package handlers

import (
	"fmt"
	"log"
	"my-drive/internal/app/db"
	"my-drive/internal/pkg/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Directory struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type DirDel struct {
	Node lib.FsNode `json:"node"`
	Soft bool       `json:"soft"`
}

func CreateDirectory(ctx *gin.Context) {
	var json lib.FsNode

	claims, err := lib.RetrieveJWTClaims(ctx)
	if err != nil {
		log.Println("Error getting user claims", err)
		lib.RespondError(ctx, http.StatusForbidden, "Claims error")
		return
	}

	if err := ctx.ShouldBindJSON(&json); err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	// TODO: Make a better way to handle this error
	if err := lib.CreateDir(claims, json); err != nil {
		lib.RespondError(ctx, http.StatusInternalServerError, "ill need to handle it better")
		return
	}

	lib.Responder(ctx, http.StatusOK, "Create Directory Successfully", nil)
}

func FindDirectory(ctx *gin.Context) {
	database := db.DB.Connection
	var dir Directory
	var dirFiles []lib.FsNode

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
		var file lib.FsNode
		if err := rows.Scan(&file.Name, &file.Type, &file.Size, &file.Created_at, &file.Updated_at); err != nil {
			fmt.Printf("%s", err)
			lib.RespondError(ctx, http.StatusInternalServerError, "Error reading files")
			return
		}

		dirFiles = append(dirFiles, file)
	}
}

func DeleteDirectory(ctx *gin.Context) {
	var json DirDel
	var dir Directory

	if err := ctx.ShouldBindJSON(&json); err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	if err := ctx.ShouldBindUri(&dir); err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	json.Node.ID = &dir.ID

	claims, err := lib.RetrieveJWTClaims(ctx)
	if err != nil {
		log.Println("Error getting user claims", err)
		lib.RespondError(ctx, http.StatusForbidden, "Claims error")
		return
	}

	if json.Soft {
		if deleteErr := lib.SoftDeleteDir(claims, json.Node); deleteErr != nil {
			lib.RespondError(ctx, http.StatusInternalServerError, "Error deleting directory")
			return
		}
	} else {
		if deleteErr := lib.DeleteDir(claims, json.Node); deleteErr != nil {
			lib.RespondError(ctx, http.StatusInternalServerError, "Error deleting directory")
			return
		}
	}

	lib.Responder(ctx, http.StatusOK, "Deleted Directory Successfully", nil)
}
