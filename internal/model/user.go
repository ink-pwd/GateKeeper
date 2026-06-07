package model

type User struct {
	//минимальная валидация получаемых данных
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
