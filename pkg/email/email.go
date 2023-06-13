package email

import (
	"fmt"
	"io"
	"net/smtp"
	"strconv"
	"strings"
)

type Email struct {
	From          string
	To            []string
	Subject       string
	PlainTextBody string
}

type Emailer func(*Email) error

func NewMockEmailer(w io.Writer, injectErr error) Emailer {
	return func(email *Email) error {
		msg := fmt.Sprintf("New email: \n\tFrom: %s\n\tTo: %s\n\tSubject: %s\n\tBody:\n\n%s\n",
			email.From,
			email.To,
			email.Subject,
			email.PlainTextBody,
		)
		if injectErr != nil {
			return injectErr
		}
		_, err := w.Write([]byte(msg))
		return err
	}
}

type SMTPEmailerConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Sender   string `json:"sender"`
	Password string `json:"password"`
}

func NewSMTPEmailer(config *SMTPEmailerConfig) Emailer {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	return func(email *Email) error {
		if email.From == "" {
			email.From = config.Sender
		}
		addr := config.Host + ":" + strconv.Itoa(config.Port)
		return smtp.SendMail(addr, auth, config.Username, email.To, []byte(SMTPMsg(email)))
	}
}

// generates the message string that will be sent to the SMTP server
func SMTPMsg(e *Email) string {
	headerMap := map[string]string{
		"From":         e.From,
		"To":           strings.Join(e.To, "; "),
		"Subject":      e.Subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain",
	}
	header := ""
	for key, val := range headerMap {
		header += key + ":" + val + "\r\n"
	}
	body := e.PlainTextBody
	return header + "\r\n" + body + "\r\n"
}
