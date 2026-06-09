package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(email, secret string, timeDuration int) (string, error) {
	var (
		token *jwt.Token
	)
	/*
		Передаем метод + параметры, которые будут вшиты в jwt
		Создаем токен на 24 часа(exp)
	*/
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Duration(timeDuration) * time.Hour).Unix(),
	})

	/*
		Подписываем нашим ключом
	*/
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString, secret string) (string, error) {
	var (
		token  *jwt.Token
		err    error
		claims jwt.MapClaims
		email  string
	)
	/*
		Парсим токен, проверяя подпись секрет и exp
	*/
	token, err = jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	/*
		Проверяем валидность токена
	*/
	if err != nil || !token.Valid {
		return "", err
	}

	/*
		Возвращаем email если валидно
	*/
	claims = token.Claims.(jwt.MapClaims)
	email = claims["email"].(string)
	return email, nil
}
