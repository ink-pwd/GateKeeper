package model

type User struct {
	/*
		Минимальная валидация получаемых данных
	*/
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
