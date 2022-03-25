package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func Sign(secretKey string, payload map[string]interface{}, exp int64) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": exp,
	}

	for k, v := range payload {
		claims[k] = v
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
}

func Parse(secretKey, tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(jwt *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, nil
}
