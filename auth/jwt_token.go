package auth

import (
	"time"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

const jwt_secret = "secret"

// only accessible through login
func GenToken(name string) (string, error) {
	claims := Claims{
		name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwt_secret))
	return tokenString, err
}

// verify token and return username if valid
func ValidateToken(tokenString string) (string, error) {
	claims := &Claims{}
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", jwt.ErrInvalidKeyType
		}
		return []byte(jwt_secret), nil
	}
	token, err := jwt.ParseWithClaims(tokenString, claims, keyfunc)
	if err != nil || !token.Valid {
		return "", err
	}

	log.Printf("Username: %s has been verified!", claims.Username)

	return claims.Username, nil
}
