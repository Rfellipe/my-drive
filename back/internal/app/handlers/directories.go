package handlers

import (
	"fmt"
	"log"
	"my-drive/internal/app/models"
	"my-drive/internal/pkg/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Directory struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func CreateDirectory(ctx *gin.Context) {
	var json models.FsNode

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

func ListDirectoryFile(ctx *gin.Context) {
	var dir Directory

	if err := ctx.ShouldBindUri(&dir); err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	claims, error := lib.RetrieveJWTClaims(ctx)
	if error != nil {
		lib.RespondError(ctx, http.StatusForbidden, "No user ID found")
		return
	}

	nodes, error := lib.ListFiles(claims, dir.ID)
	if error != nil {
		fmt.Println(error)
		lib.RespondError(ctx, http.StatusInternalServerError, "Error getting files for this directory")
		return
	}

	lib.Responder(ctx, http.StatusOK, "Files found", nodes)
}

func UpdateDirectory(ctx *gin.Context) {
	var json models.DirUpdate
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

	if updateErr := lib.Rename(claims, json.Node, json.OldPath, json.NewPath); updateErr != nil {
		lib.RespondError(ctx, http.StatusInternalServerError, "Error updating directory")
		return
	}

	lib.Responder(ctx, http.StatusOK, "Updated Directory Successfully", nil)
}

func DeleteDirectory(ctx *gin.Context) {
	var json models.DirDel
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
