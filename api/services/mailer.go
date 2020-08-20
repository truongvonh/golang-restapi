package services

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"os"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailData struct {
	UserName string
	UserMail string
	Content  string
}

func SendMail(mailParam MailData) (*rest.Response, error) {
	from := mail.NewEmail(os.Getenv("MAIL_USER"), os.Getenv("MAIL_DEFAULT"))
	subject := "Confirm your account with verify code"
	to := mail.NewEmail(mailParam.UserName, mailParam.UserMail)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>" + mailParam.Content + "</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("MAIL_KEY"))
	return client.Send(message)

}
