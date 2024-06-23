package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

const jwt_secret = "secret"

func Gen_token(name string) (string, error) {
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

func Validate_token(tokenString string) (string, error) {
	claims := &Claims{}
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", jwt.ErrInvalidKeyType
		}
		return []byte(jwt_secret), nil
	}
	token, err := jwt.ParseWithClaims(tokenString, claims, keyfunc)
	if err != nil || !token.Valid{
		return "", err
	}

	// remove in production
	fmt.Println(fmt.Sprintf("Username: %s has been verified!", claims.Username))

	return claims.Username, nil
}
