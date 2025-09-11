package handlers

import (
	"fmt"
	_ "fmt"
	"my-drive/internal/app/db"
	"my-drive/internal/pkg/lib"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserForm struct {
	Email string `json:"email" binding:"required"`
	Pass  string `json:"pass" binding:"required"`
}

func RegistrationHandler(ctx *gin.Context) {
	database := db.DB.Connection
	var json UserForm

	if err := ctx.Bind(&json); err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, "Invalid Input")
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(json.Pass), bcrypt.DefaultCost)
	if err != nil {
		lib.RespondError(ctx, http.StatusInternalServerError, "Hash error, account not created")
		return
	}

	userid := ""
	err = database.QueryRow(`INSERT INTO users(email, password) VALUES($1, $2) RETURNING id`, json.Email, hashedPass).Scan(&userid) 
	if err != nil {
		fmt.Printf("%s", err)
		lib.RespondError(ctx, http.StatusInternalServerError, "Insert error, account not created")
		return
	}

	_, err = database.Exec(`INSERT INTO nodes(name, type, owner) VALUES ($1, $2, $3)`, userid, "dir", userid)
	if err != nil {
		fmt.Printf("%s", err)
		lib.RespondError(ctx, http.StatusInternalServerError, "Insert error, your dir was not created")
		database.Exec(`DELETER FROM users WHERE id = $1`, userid)
		return
	}

	err = os.Mkdir(fmt.Sprintf("%s/%s/", os.Getenv("DRIVE_ROOT"), userid), os.ModePerm)
	if err != nil {
		fmt.Printf("%s", err)
		lib.RespondError(ctx, http.StatusInternalServerError, "your dir was not created")
		return
	}

	lib.Responder(ctx, http.StatusOK, "Account created", nil)
}

func LoginHandler(ctx *gin.Context) {
	database := db.DB.Connection
	var userInfo lib.UserSubject
	var json UserForm

	if err := ctx.Bind(&json); err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, "Invalid Input")
		return
	}

	hashedPassword := "$2a$12$VSZFnXJazJ.2WE1hVrydMu6sYUPWQLfdlZs5ivFVPzk4BDVTOLQAy" // Dummy bcrypt hash
	userExists := true

	err := database.QueryRow(
		`SELECT id, email, password, status, login_attempts FROM users WHERE email=$1`,
		json.Email,
	).Scan(
		&userInfo.Id,
		&userInfo.Email,
		&hashedPassword,
		&userInfo.Status,
		&userInfo.Login_attempts,
	)
	if err != nil {
		userExists = false
	}

	// Dummy hashing to avoid timming attacks
	_ = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(json.Pass))

	if !userExists {
		lib.Responder(ctx, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	if userInfo.Status == "banned" || userInfo.Login_attempts >= 5 {
		lib.Responder(ctx, http.StatusUnauthorized, "Account banned or locked", nil)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(json.Pass)) != nil {
		_, _ = database.Exec(`UPDATE users SET login_attempts=$1 WHERE id=$2`, userInfo.Login_attempts+1, userInfo.Id)
		lib.Responder(ctx, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	database.Exec(`UPDATE users SET login_attempts=$1 WHERE id=$2`, 0, userInfo.Id)

	token, _ := lib.GenerateJWT(userInfo)
	ctx.Header("access_token", token)
	// ctx.Header("expires", fmt.Sprintf("%d", exp))
	lib.Responder(ctx, http.StatusOK, "Login successful", nil)
}



