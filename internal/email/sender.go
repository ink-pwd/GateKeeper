package email

import (
	"fmt"
	"net/smtp"
)

type Sender struct {
	email    string
	password string
	host     string
	port     string
}

func NewSender(email, password, host, port string) *Sender {
	return &Sender{
		email:    email,
		password: password,
		host:     host,
		port:     port,
	}
}

func (s *Sender) SendMessage(email, link string, time int) error {
	var (
		auth smtp.Auth
		err  error
		addr []string
	)
	//формируем массив адресов для отправки(в нашем случае один)
	addr = append(addr, email)
	auth = smtp.PlainAuth(
		"",
		s.email,
		s.password,
		s.host,
	)

	//отправляем письмо
	err = smtp.SendMail(
		s.host+s.port,
		auth,
		s.email,
		addr,
		[]byte(fmt.Sprintf("Subject\nVerification link: %s\nWarning: the link is valid for %d minutes", link, time)),
	)
	return err
}
