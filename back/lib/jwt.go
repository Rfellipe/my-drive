package lib

import (
	"my-drive/dtos"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userInfo dtos.UserInfo) (string, int64) {
	var (
		key []byte
		tk  *jwt.Token
		s   string
	)

	now := time.Now().UnixMilli()
	later := time.Now().Add(1 * time.Hour).UnixMilli()

	key = []byte(os.Getenv("JWT_SECRET"))
	tk = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iat": now,
			"exp": later,
			"userInfo": userInfo,
		},
	)
	s, _ = tk.SignedString(key)
	
	return s, later
}
