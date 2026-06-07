package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/ink-pwd/Gatekeeper/internal/consts"
	"github.com/ink-pwd/Gatekeeper/internal/handler"
	"github.com/ink-pwd/Gatekeeper/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	var (
		err     error
		mux     *http.ServeMux
		log     logger.Logger
		handl   handler.AuthHandler
		server  string
		connStr string
		secret  string
		db      *sql.DB
	)

	//инициализация логера
	log = logger.NewStdLogger()

	//загружаем окружение
	err = godotenv.Load()
	if err != nil {
		log.Fatal("environment startup error: ", err.Error())
		return
	}
	server = os.Getenv("SERVER")
	connStr = os.Getenv("DB")
	secret = os.Getenv("SECRET")

	//подключаемся к бд
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("error db connect: ", err.Error())
		return
	}
	log.Info("success connect to db")

	//инициализация зависимостей
	mux = http.NewServeMux()
	handl = *handler.NewAuthHandler(log, db, secret)

	//регистрируем маршруты
	mux.HandleFunc(consts.REGISTER, handl.Register)
	mux.HandleFunc(consts.LOGIN, handl.Login)
	mux.HandleFunc(consts.VALIDATE, handl.Validate)
	//запуск сервера
	err = http.ListenAndServe(server, mux)
	if err != nil {
		log.Fatal("the server did not start: ", err.Error())
		return
	}
	log.Info("server has been started at ", server)
}
