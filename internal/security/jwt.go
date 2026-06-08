package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(email string, secret string) (string, error) {
	var (
		token *jwt.Token
	)
	//передаем метод + параметры, которые будут вшиты в jwt
	//создаем токен на 24 часа(exp)
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // +24 часа
	})

	//подписываем нашим ключом
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string, secret string) (string, error) {
	var (
		token  *jwt.Token
		err    error
		claims jwt.MapClaims
		email  string
	)
	//парсим токен, проверяя подпись секрет и exp
	token, err = jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	//проверяем валидность токена
	if err != nil || !token.Valid {
		return "", err
	}

	//возвращаем email если валидно
	claims = token.Claims.(jwt.MapClaims)
	email = claims["email"].(string)
	return email, nil
}
