package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ink-pwd/Gatekeeper/internal/consts"
	"github.com/ink-pwd/Gatekeeper/internal/email"
	"github.com/ink-pwd/Gatekeeper/internal/model"
	"github.com/ink-pwd/Gatekeeper/internal/security"
	"github.com/ink-pwd/Gatekeeper/internal/storage"
	"github.com/ink-pwd/Gatekeeper/logger"
)

type AuthHandler struct {
	log          logger.Logger
	db           *sql.DB
	secret       string
	clientRedis  *storage.ClientRedis
	sender       *email.Sender
	timeDuration int
	protocol     string
	server       string
	port         string
}

func NewAuthHandler(log logger.Logger, db *sql.DB, clientRedis *storage.ClientRedis,
	sender *email.Sender, timeDuration int, protocol, server, secret, port string) *AuthHandler {
	return &AuthHandler{
		log:          log,
		db:           db,
		secret:       secret,
		clientRedis:  clientRedis,
		sender:       sender,
		timeDuration: timeDuration,
		protocol:     protocol,
		server:       server,
		port:         port,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		body  []byte
		uuid  string
		exist int64
		user  *model.User
		link  string
	)
	/*получаем json запроса, его будем использовать для записи в бд
	в случае подтверждения почты, а как временное хранилище данных используем
	Redis, в качестве ключа будет использоваться uuid*/

	//получаем body запроса
	err = json.NewDecoder(r.Body).Decode(&user)

	//закрываем чтение тела запроса
	r.Body.Close()
	if err != nil {
		h.log.Error("json decode error: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//сразу хешируем пароль пользователя!
	user.Password, err = security.CryptoHash(user.Password)
	if err != nil {
		h.log.Error("crypho hash err: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err = json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for {
		uuid = security.GetUUID()
		exist, err = h.clientRedis.Exist(uuid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if exist == 0 {
			break
		}

	}
	/*записываем в redis
	uuid => body request json
	timeDuration - время на подтверждение почты*/
	h.clientRedis.Set(uuid, body, h.timeDuration)

	/*Асинхронно отправляем письмо,
	что бы пользователь долго не ждал ответ от свервера*/
	link = fmt.Sprintf("%s://%s%s%s?token=%s", h.protocol, h.server, h.port, consts.VERIFY, uuid)
	go h.sender.SendMessage(user.Email, link, h.timeDuration)

	h.log.Info("success send code: %s", body)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var (
		user          *model.User
		jwtToken      string
		err           error
		password_hash string
	)
	//получаем информацию от пользователя

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//получаем body запроса
	err = json.NewDecoder(r.Body).Decode(&user)
	r.Body.Close()

	if err != nil {
		h.log.Error("json decode error: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	password_hash = storage.GetUserByEmail(user.Email, user.Password, h.db, h.log)
	if !security.Compare(user.Password, password_hash) {
		/*возвращаем 401, что значит, что данные введены неверно
		никаких уточнений не делаем для защиты,
		пользователь не должен знать имеется ли такой email в базе или нет*/
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//получаем jwt-токен
	jwtToken, err = security.CreateToken(user.Email, h.secret)
	if err != nil {
		h.log.Error("token creation error: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	h.log.Info("success login user: ", user.Email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": jwtToken})
}

func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	var (
		jwtToken string
		email    string
		err      error
	)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jwtToken = r.Header.Get("Authorization")
	if jwtToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//убираем префикс "Bearer "
	if len(jwtToken) > 7 && jwtToken[:7] == "Bearer " {
		jwtToken = jwtToken[7:]
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//валидируем токен
	email, err = security.ValidateToken(jwtToken, h.secret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//токен валиден, возвращаем email
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"email": email,
		"valid": "true",
	})
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var (
		token    string
		err      error
		data     string
		userByte []byte
		user     *model.User
		ok       bool
	)

	h.log.Info("VERIFY HANDLER CALLED, URL: ", r.URL.String())
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//получаем токен
	token = r.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//получаем пользователя
	data, err = h.clientRedis.Get(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//получаем информацию о пользователе
	userByte = []byte(data)
	err = json.Unmarshal(userByte, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error("json: ", err.Error())
		return
	}

	//удаляем пользователя
	err = h.clientRedis.Del(token)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error("redis del: ", err.Error())
		return
	}

	//записываем информацию в бд
	ok = storage.AddUser(user.Email, user.Password, h.db, h.log)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
