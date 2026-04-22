package host

import (
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"nlhosting/cfg"
)

func SendMail(receiver, title, body string) {
	to := []string{receiver}
	subject := fmt.Sprintf("NodeLoc Free Hosting %s", title)

	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		cfg.Config.Mail.User, to[0], subject, body,
	)

	auth := smtp.PlainAuth(
		"",
		cfg.Config.Mail.User,
		cfg.Config.Mail.Pass,
		cfg.Config.Mail.Host,
	)

	addr := net.JoinHostPort(
		cfg.Config.Mail.Host,
		cfg.Config.Mail.Port,
	)

	err := smtp.SendMail(addr, auth, cfg.Config.Mail.User, to, []byte(message))
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
