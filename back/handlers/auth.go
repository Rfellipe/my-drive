package handlers

import (
	"fmt"
	"my-drive/dtos"
	"my-drive/lib"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegistrationHandler(ctx *gin.Context) {
	db := lib.DB.Connection
	var json dtos.UserForm

	if err := ctx.Bind(&json); err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, "Invalid Input")
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(json.Pass), bcrypt.DefaultCost)
	if err != nil {
		lib.RespondError(ctx, http.StatusInternalServerError, "Hash error, account not created")
		return
	}

	_, err = db.Exec(`INSERT INTO users(email, password) VALUES($1, $2)`, json.Email, hashedPass)
	if err != nil {
		lib.RespondError(ctx, http.StatusInternalServerError, "Insert error, account not created")
		return
	}

	lib.Responder(ctx, http.StatusOK, "Account created", nil)
}

func LoginHandler(ctx *gin.Context) {
	db := lib.DB.Connection
	var userInfo dtos.UserInfo
	var json dtos.UserForm

	if err := ctx.Bind(&json); err != nil {
		lib.RespondError(ctx, http.StatusBadRequest, "Invalid Input")
		return
	}

	hashedPassword := "$2a$12$VSZFnXJazJ.2WE1hVrydMu6sYUPWQLfdlZs5ivFVPzk4BDVTOLQAy" // Dummy bcrypt hash
	userExists := true

	err := db.QueryRow(
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
		_, _ = db.Exec(`UPDATE users SET login_attempts=$1 WHERE id=$2`, userInfo.Login_attempts+1, userInfo.Id)
		lib.Responder(ctx, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	db.Exec(`UPDATE users SET login_attempts=$1 WHERE id=$2`, 0, userInfo.Id)

	token, exp := lib.GenerateJWT(userInfo)
	ctx.Header("access_token", token)
	ctx.Header("expires", fmt.Sprintf("%d", exp))
	lib.Responder(ctx, http.StatusOK, "Login successful", nil)
}
