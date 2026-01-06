package lib

import (
	"fmt"
	"log"
	"my-drive/internal/app/db"
	"os"
	"strings"

	_ "github.com/gin-gonic/gin"
)

func GetUserDir(userId string) string {
	var root string = os.Getenv("DRIVE_ROOT")
	var userRootDir string = fmt.Sprintf("%s/%s", root, userId)

	return userRootDir
}

func CreateDir(userClaims UserSubject, dirPath string, parentId string) error {
	database := db.DB.Connection
	var fullPath string = fmt.Sprintf("%s/%s", GetUserDir(userClaims.Id), dirPath)
	var names []string = strings.Split(dirPath, "/")

	if parentId != "" {
		_, err := database.Exec(
			`INSERT INTO
				nodes(name, type, owner, parent_id)
			VALUES ($1, $2, $3, $4)`,
			names[len(names)-1], Directory.String(),
			userClaims.Id, parentId,
		)
		if err != nil {
			log.Printf("Insert error, failed to create directory: %s", err)
			return err
		}
	} else {
		_, err := database.Exec(
			`INSERT INTO
				nodes(name, type, owner, parent_id)
			VALUES ($1, $2, $3, $4)`,
			names[len(names)-1], Directory.String(),
			userClaims.Id, userClaims.RootDirId,
		)
		if err != nil {
			log.Printf("Insert error, failed to create 1directory: %s", err)
			return err
		}
	}

	err := os.Mkdir(fullPath, os.ModePerm)
	if err != nil {
		log.Printf("Insert error, failed to create directory: %s", err)
		return err
	}

	return nil
}

func DeleteDir(dirPath string) error {

	return nil
}
