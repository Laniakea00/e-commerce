package utils

import (
	"fmt"
	"net/smtp"
)

var (
	SMTPServer = "smtp.gmail.com"
	SMTPPort   = 587
	EmailUser  = "alan06b@gmail.com"
	EmailPass  = "fxznbzbsskixkiwz"
)

// SendEmailVerification отправляет письмо с подтверждением
func SendEmailVerification(to string, verificationLink string) error {
	auth := smtp.PlainAuth("", EmailUser, EmailPass, SMTPServer)

	subject := "Confirm your registration"
	body := fmt.Sprintf(`
Hello,

Thank you for registering. Please confirm your email by clicking the link below:

%s

If you did not register, please ignore this message.
`, verificationLink)

	msg := "From: " + EmailUser + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	addr := fmt.Sprintf("%s:%d", SMTPServer, SMTPPort)
	return smtp.SendMail(addr, auth, EmailUser, []string{to}, []byte(msg))
}
