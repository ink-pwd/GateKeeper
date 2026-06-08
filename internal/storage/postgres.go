package storage

import (
	"database/sql"

	"github.com/ink-pwd/Gatekeeper/logger"
)

func AddUser(email, password_hash string, db *sql.DB, log logger.Logger) bool {
	var (
		err error
	)

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

func GetUserByEmail(email, password string, db *sql.DB, log logger.Logger) string {
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
		return ""
	}

	return password_hash

}
