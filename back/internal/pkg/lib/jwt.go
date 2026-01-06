package lib

import (
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserSubject struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	Status         string `json:"status"`
	RootDirId      string `json:"rootDirId"`
	Login_attempts int    `json:"loginAttempts"`
}

func GenerateJWT(userInfo UserSubject) (string, int64) {
	var (
		key []byte
		tk  *jwt.Token
		s   string
	)

	now := time.Now().Unix()
	later := time.Now().Add(time.Hour * 24).Unix() // 24 hours

	key = []byte(os.Getenv("JWT_SECRET"))
	tk = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iat": now,
			"exp": later,
			"sub": userInfo,
		},
	)
	s, _ = tk.SignedString(key)

	return s, later
}

func ValidateJWT(s string) (*jwt.Token, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	return jwt.Parse(s, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return key, nil
	})
}

func RetrieveJWTClaims(ctx *gin.Context) (UserSubject, error) {
	var claims UserSubject

	id, _ := ctx.Get("userId")
	email, _ := ctx.Get("userEmail")
	status, _ := ctx.Get("userStatus")
	rootDirId, _ := ctx.Get("userRootDir")
	loginAttempts, _ := ctx.Get("userLoginAttempts")

	if id == nil ||
		email == nil ||
		status == nil ||
		rootDirId == nil ||
		loginAttempts == nil {
		return claims, errors.New("Missing info on token")
	}

	claims.Id = id.(string)
	claims.Email = email.(string)
	claims.Status = status.(string)
	claims.RootDirId = rootDirId.(string)
	claims.Login_attempts = loginAttempts.(int)

	return claims, nil
}
