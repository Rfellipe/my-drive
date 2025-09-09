package lib

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserSubject struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	Status         string `json:"status"`
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
			// "id": userInfo.Id,
			// "email:": userInfo.Email,
			// "status": userInfo.Status,
			// "login_attempts": userInfo.Login_attempts,
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
