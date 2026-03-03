package email

import (
	"fmt"
	"net/smtp"
	"strconv"
)

type Config struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}

type Sender struct {
	cfg Config
}

func New(cfg Config) *Sender {
	return &Sender{cfg: cfg}
}

func (s *Sender) SendResetCode(to, code string) error {
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.SMTPHost)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Password Reset Code\r\n\r\nYour password reset code is: %s\r\nThis code expires in 15 minutes.",
		s.cfg.From, to, code,
	)

	addr := s.cfg.SMTPHost + ":" + strconv.Itoa(s.cfg.SMTPPort)
	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, []byte(msg))
}
