package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ink-pwd/Gatekeeper/internal/model"
	"github.com/ink-pwd/Gatekeeper/internal/storage"
	"github.com/ink-pwd/Gatekeeper/internal/token"
	"github.com/ink-pwd/Gatekeeper/logger"
)

type AuthHandler struct {
	Log    logger.Logger
	DB     *sql.DB
	Secret string
}

func NewAuthHandler(log logger.Logger, db *sql.DB, secret string) *AuthHandler {
	return &AuthHandler{
		Log:    log,
		DB:     db,
		Secret: secret,
	}
}

func (h *AuthHandler) parseUser(w http.ResponseWriter, r *http.Request) (*model.User, bool) {
	var (
		user model.User
		err  error
	)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return nil, false
	}

	//получаем body запроса
	err = json.NewDecoder(r.Body).Decode(&user)
	r.Body.Close()

	if err != nil {
		h.Log.Error("json decode error: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return nil, false
	}

	/*минимальная валидация email
	убираем комменты только если удалили валидацию в user.go struct
	иначе она просто не имеет смысла

	if !strings.Contains(user.Email, "@") || !strings.Contains(user.Email, ".") {
		h.Log.Error("invalid email: ", user.Email)
		w.WriteHeader(http.StatusBadRequest)
		return nil, false
	}*/

	return &user, true
}
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var (
		user *model.User
		ok   bool
	)

	//получаем информацию от пользователя
	user, ok = h.parseUser(w, r)
	if !ok {
		return
	}

	//сохраняем в бд
	if !storage.AddUser(user.Email, user.Password, h.DB, h.Log) {
		w.WriteHeader(http.StatusConflict)
		return
	}

	h.Log.Info("success add user: %s", user.Email)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var (
		user     *model.User
		ok       bool
		jwtToken string
		err      error
	)
	//получаем информацию от пользователя
	user, ok = h.parseUser(w, r)
	if !ok {
		return
	}

	if !storage.GetUserByEmail(user.Email, user.Password, h.DB, h.Log) {
		/*возвращаем 401, что значит, что данные введены неверно
		никаких уточнений не делаем для защиты,
		пользователь не должен знать имеется ли такой email в базе или нет*/
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//получаем jwt-токен
	jwtToken, err = token.CreateToken(user.Email, h.Secret)
	if err != nil {
		h.Log.Error("token creation error: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	h.Log.Info("success login user: ", user.Email)
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
	email, err = token.ValidateToken(jwtToken, h.Secret)
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
