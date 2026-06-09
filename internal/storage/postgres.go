package storage

import (
	"database/sql"

	"github.com/ink-pwd/Gatekeeper/logger"
)

func AddUser(email, password_hash string, db *sql.DB, log logger.Logger) bool {
	var (
		err error
	)

	/*
		Делаем запись в бд
	*/
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
	/*
		Ищем пароль пользователя по email
		проверяя верифицирован ли он
	*/
	err = db.QueryRow(
		`SELECT password_hash FROM public."User" WHERE email = $1 AND verify = true`,
		email,
	).Scan(&password_hash)

	if err != nil {
		return ""
	}

	return password_hash

}

func CheckUser(email string, db *sql.DB, log logger.Logger) bool {
	var (
		exists bool
		err    error
	)
	/*
		Проверяем есть ли пользователь
	*/
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM public."User" WHERE email = $1)`,
		email).Scan(&exists)

	if err != nil {
		log.Error("error checking user existence: %s", err.Error())
		return true
	}

	return exists
}

func VerifyUser(email string, db *sql.DB, log logger.Logger) {
	var (
		err error
	)
	/*
		Меняем значение верификации на тру
	*/
	_, err = db.Exec(`UPDATE public."User" SET verify = true WHERE email = $1`, email)
	if err != nil {
		log.Error("failed to update user status for %s: %s", email, err.Error())
	}
}

func DeleteNonVerify(db *sql.DB, log logger.Logger) {
	var (
		err error
	)
	/*
		Удаляем всех не верифицированных пользоваетелей
	*/
	_, err = db.Exec(`DELETE public."USER" WHERE verify == false`)
	if err != nil {
		log.Error("error delete user: ", err.Error())
	}
}
