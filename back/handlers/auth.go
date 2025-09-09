package handlers

import (
	_ "fmt"
	"my-drive/lib"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserForm struct {
	Email string `json:"email" binding:"required"`
	Pass  string `json:"pass" binding:"required"`
}

func RegistrationHandler(ctx *gin.Context) {
	db := lib.DB.Connection
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

	_, err = db.Exec(`INSERT INTO users(email, password) VALUES($1, $2)`, json.Email, hashedPass)
	if err != nil {
		lib.RespondError(ctx, http.StatusInternalServerError, "Insert error, account not created")
		return
	}

	lib.Responder(ctx, http.StatusOK, "Account created", nil)
}

func LoginHandler(ctx *gin.Context) {
	db := lib.DB.Connection
	var userInfo lib.UserInfo
	var json UserForm

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

	token, _ := lib.GenerateJWT(userInfo)
	ctx.Header("access_token", token)
	// ctx.Header("expires", fmt.Sprintf("%d", exp))
	lib.Responder(ctx, http.StatusOK, "Login successful", nil)
}

func JWTAuthorizeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			lib.RespondError(ctx, http.StatusUnauthorized, "Missing Authorization header")
			ctx.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			lib.RespondError(ctx, http.StatusUnauthorized, "Invalid Authorization header format")
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		claims, err := lib.ValidateJWT(tokenString)
		if err != nil {
			lib.RespondError(ctx, http.StatusUnauthorized, "Invalid token")
			ctx.Abort()
			return
		}

		exp, err := claims.Claims.GetExpirationTime()
		if err != nil || exp == nil || exp.Time.Before(time.Now()) {
			lib.RespondError(ctx, http.StatusUnauthorized, "Token expired")
			ctx.Abort()
			return
		}

		sub, ok := claims.Claims.(jwt.MapClaims)["sub"].(map[string]interface{})
		if !ok {
			lib.RespondError(ctx, http.StatusUnauthorized, "Invalid subject claim")
			ctx.Abort()
			return
		}

		userID, _ := sub["id"].(string)
		email, _ := sub["email"].(string)
		status, _ := sub["status"].(string)
		loginAttempts, _ := sub["loginAttempts"].(float64) 

		ctx.Set("userId", userID)
		ctx.Set("userEmail", email)
		ctx.Set("userStatus", status)
		ctx.Set("userLoginAttempts", int(loginAttempts))
		ctx.Set("claims", claims.Claims) 

		ctx.Next()
	}
}

