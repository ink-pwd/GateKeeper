package storage

import (
	"database/sql"

	"github.com/ink-pwd/Gatekeeper/logger"
	"golang.org/x/crypto/bcrypt"
)

func AddUser(email, password string, db *sql.DB, log logger.Logger) bool {
	var (
		password_hash string
		err           error
	)

	//превращаем пароль в хеш
	password_hash, err = cryptoHash(password)
	if err != nil {
		log.Error("failed to create password hash: ", err.Error())
		return false
	}

	//делаем запись в бд
	_, err = db.Exec(
		`INSERT INTO public."User" (email, password_hash) VALUES ($1, $2)`,
		email,
		password_hash,
	)
	if err != nil {
		log.Error("failed add user: ", err.Error())
		return false
	}
	return true
}

func GetUserByEmail(email, password string, db *sql.DB, log logger.Logger) bool {
	var (
		password_hash string
		err           error
	)
	//получаем пользователя по email и получаем хеш-пароль
	err = db.QueryRow(
		`SELECT password_hash FROM public."User" WHERE email = $1`,
		email,
	).Scan(&password_hash)

	if err != nil {
		log.Error("user not found: %s", err.Error())
		return false
	}

	//сравниваем наш пароль с хешем
	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	if err != nil {
		log.Error("invalid password for %s", email)
		return false
	}

	return true

}

func cryptoHash(password string) (string, error) {
	var (
		password_hash []byte
		err           error
	)
	password_hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(password_hash), nil
}
