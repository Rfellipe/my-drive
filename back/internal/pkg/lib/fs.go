package lib

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"my-drive/internal/app/db"
	"my-drive/internal/app/models"
	"os"
	"time"

	_ "github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type FsTypes int

const (
	Directory FsTypes = iota
	File
)

var FsTypeName = map[FsTypes]string{
	Directory: "dir",
	File:      "file",
}

type paths struct {
	old string
	new string
}

func (ss FsTypes) String() string {
	return FsTypeName[ss]
}

func getUserDir(userId string) string {
	var root string = os.Getenv("DRIVE_ROOT")
	var userRootDir string = fmt.Sprintf("%s/%s", root, userId)

	return userRootDir
}

func moveFileToBin(userId string, filePath string) error {
	var fullPath string = fmt.Sprintf("%s/%s/", getUserDir(userId), filePath)
	var binPath string = fmt.Sprintf("%s/%s/%s", getUserDir(userId), "recycle_bin", filePath)

	if err := os.Rename(fullPath, binPath); err != nil {
		return err
	}

	return nil
}

func CreateUserDir(userId string) error {
	var userDir string = getUserDir(userId)

	err := os.MkdirAll(fmt.Sprintf("%s/%s/", userDir, "recycle_bin"), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func CreateDir(userClaims UserSubject, node models.FsNode) error {
	database := db.DB.Connection
	var fullPath string = fmt.Sprintf("%s/%s", getUserDir(userClaims.Id), node.Name)

	if node.ParentId != nil {
		_, err := database.Exec(
			`INSERT INTO
				nodes(name, type, owner, parent_id)
			VALUES ($1, $2, $3, $4)`,
			node.Name, Directory.String(),
			userClaims.Id, node.ParentId,
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
			node.Name, Directory.String(),
			userClaims.Id, userClaims.RootDirId,
		)
		if err != nil {
			log.Printf("Insert error, failed to create 1directory: %s", err)
			return err
		}
	}

	err := os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		log.Printf("Insert error, failed to create directory: %s", err)
		return err
	}

	return nil
}

func SoftDeleteDir(userClaims UserSubject, node models.FsNode) error {
	database := db.DB.Connection
	ctx, cancel := context.WithCancel(db.DB.RootContext)

	tx, err := database.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		cancel()
		return err
	}

	rows, queryErr := tx.QueryContext(
		ctx, `SELECT id FROM nodes WHERE id = $1 OR parent_id = $1`, node.ID,
	)
	if queryErr != nil {
		tx.Rollback()
		cancel()
		return queryErr
	}
	defer rows.Close()

	var nodes []string
	for rows.Next() {
		var id string
		if rowErr := rows.Scan(&id); rowErr != nil {
			fmt.Println(rowErr)
			cancel()
			return rowErr
		}

		nodes = append(nodes, id)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		fmt.Println(rowsErr)
		cancel()
		return rowsErr
	}

	_, execErr := tx.ExecContext(
		ctx, `UPDATE nodes SET deleted_at = $1 WHERE id = any($2)`,
		time.Now(), pq.Array(nodes),
	)
	if execErr != nil {
		fmt.Println(execErr)
		cancel()
		return execErr
	}

	moveErr := moveFileToBin(userClaims.Id, node.Name)
	if moveErr != nil {
		fmt.Println(moveErr)
		cancel()
		return moveErr
	}

	if err := tx.Commit(); err != nil {
		cancel()
		return err
	}

	cancel()
	return nil
}

func DeleteDir(userClaims UserSubject, node models.FsNode) error {
	database := db.DB.Connection
	ctx, cancel := context.WithCancel(db.DB.RootContext)
	var filePath string = fmt.Sprintf(
		"%s/%s/%s", getUserDir(userClaims.Id), "recycle_bin", node.Name,
	)

	tx, err := database.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		cancel()
		return err
	}

	_, execErr := tx.ExecContext(
		ctx, `DELETE FROM nodes WHERE id = $1 AND deleted_at IS NOT NULL`, node.ID,
	)
	if execErr != nil {
		fmt.Println(execErr)
		cancel()
		return execErr
	}

	remErr := os.RemoveAll(filePath)
	if remErr != nil {
		cancel()
		return err
	}

	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		cancel()
		return err
	}

	cancel()
	return nil
}

func Rename(userClaims UserSubject, node models.FsNode, oldPath string, newPath string) error {
	database := db.DB.Connection
	path := paths{
		old: fmt.Sprintf(`%s/%s`, getUserDir(userClaims.Id), oldPath),
		new: fmt.Sprintf(`%s/%s`, getUserDir(userClaims.Id), newPath),
	}
	parentId := func(parentId *string) *string {
		if parentId == nil || *parentId == "" {
			return nil
		} else {
			return parentId
		}
	}

	fmt.Printf("nodeId: %s\nnewPath: %s\nparentId: %s\n", *node.ID, newPath, parentId(node.ParentId))
	_, err := database.Exec(`SELECT update_node($1, $2, $3);`,
		*node.ID, newPath, parentId(node.ParentId),
	)
	if err != nil {
		log.Printf("Error while updating node table: %s", err)
		return err
	}

	if mvErr := os.Rename(path.old, path.new); mvErr != nil {
		log.Printf("Error updating node: %s", mvErr)
		return err
	}

	return nil
}

func ListFiles(userClaims UserSubject, dirId string) ([]models.FsNode, error) {
	database := db.DB.Connection

	rows, queryErr := database.Query(`
		SELECT 
			id, name, type, size, parent_id, created_at, updated_at
		FROM
			nodes
		WHERE
			id = $1 OR parent_id = $1`,
		dirId,
	)
	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	var nodes []models.FsNode
	for rows.Next() {
		var node models.FsNode
		if rowErr := rows.Scan(
			&node.ID, &node.Name, &node.Type,
			&node.Size, &node.ParentId, &node.Created_at,
			&node.Updated_at,
		); rowErr != nil {
			return nil, rowErr
		}

		nodes = append(nodes, node)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return nodes, nil
}

func UploadFile() error {

	return nil
}
