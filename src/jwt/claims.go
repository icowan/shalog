package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
)

type ArithmeticCustomClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

var jwtKey = os.Getenv("JWT_KEY")

func init() {
	if jwtKey == "" {
		jwtKey = "hello@world:lattecake.com"
	}
}

func JwtKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(GetJwtKey()), nil
}

func GetJwtKey() string {
	return jwtKey
}
