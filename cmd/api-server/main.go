package main

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/ink-pwd/Gatekeeper/internal/consts"
	"github.com/ink-pwd/Gatekeeper/internal/email"
	"github.com/ink-pwd/Gatekeeper/internal/handler"
	"github.com/ink-pwd/Gatekeeper/internal/storage"
	"github.com/ink-pwd/Gatekeeper/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	var (
		err          error
		mux          *http.ServeMux
		log          logger.Logger
		handl        *handler.AuthHandler
		server       string
		port         string
		protocol     string
		redisPort    string
		connStr      string
		secret       string
		fromEmail    string
		password     string
		hostsmtp     string
		portsmtp     string
		time         string
		timeDuration int
		sender       *email.Sender
		db           *sql.DB
		redClient    *storage.ClientRedis
		connect      bool
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
	protocol = os.Getenv("PROTOCOL")
	port = os.Getenv("PORT")
	redisPort = os.Getenv("REDIS")
	hostsmtp = os.Getenv("HOSTSMTP")
	portsmtp = os.Getenv("PORTSMTP")
	fromEmail = os.Getenv("EMAIL")
	password = os.Getenv("PASSWORD")
	time = os.Getenv("TIMEDURATION")

	//превращаем время взятое из окружение в int
	timeDuration, err = strconv.Atoi(time)
	if err != nil {
		log.Fatal("specify an integer time duration in .env \n", err.Error())
	}

	//создаем структуру для отправки сообщений
	sender = email.NewSender(fromEmail, password, hostsmtp, portsmtp)

	//подключаемся к бд
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("error db connect: ", err.Error())
		return
	}
	log.Info("success connect to db")

	//подключаемся к redis
	redClient = storage.NewClient(redis.NewClient(&redis.Options{
		Addr:     server + redisPort,
		Password: "",
		DB:       0,
	}))

	connect = redClient.Ping()
	if !connect {
		log.Fatal("error redis connect")
	}
	log.Info("success redis connect %s%s", server, redisPort)

	//инициализация зависимостей
	mux = http.NewServeMux()
	handl = handler.NewAuthHandler(log, db, redClient, sender, timeDuration, protocol, server, secret, port)

	//регистрируем маршруты
	mux.HandleFunc(consts.REGISTER, handl.Register)
	mux.HandleFunc(consts.LOGIN, handl.Login)
	mux.HandleFunc(consts.VALIDATE, handl.Validate)
	mux.HandleFunc(consts.VERIFY, handl.Verify)
	//запуск сервера
	err = http.ListenAndServe(server+port, mux)
	if err != nil {
		log.Fatal("the server did not start: ", err.Error())
		return
	}
	log.Info("server has been started at ", server+port)
}
