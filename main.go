package handler

import (
	"net/smtp"
	"os"

	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

type Request struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Html    string `json:"html"`
}

type Response struct {
	Success bool `json:"success"`
}

type Env struct {
}

func Main(Context openruntimes.Context) openruntimes.Response {
	if Context.Req.Path == "/send" && Context.Req.Method == "POST" {
		var request Request
		err := Context.Req.BodyJson(&request)
		if err != nil {
			return Context.Res.Text("invalid body: "+err.Error(), Context.Res.WithStatusCode(500))
		}

		// Send an email via smtp
		smtpHost := os.Getenv("SMTP_HOST")
		smtpPort := os.Getenv("SMTP_PORT")
		smtpUser := os.Getenv("SMTP_USER")
		smtpPass := os.Getenv("SMTP_PASS")

		if smtpUser != "" || smtpPass != "" {
			auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
			err = smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{request.To}, []byte("Subject: "+request.Subject+"\r\n\r\n"+request.Html))
		} else {
			err = smtp.SendMail(smtpHost+":"+smtpPort, nil, "", []string{request.To}, []byte("Subject: "+request.Subject+"\r\n\r\n"+request.Html))
		}

		if err != nil {
			return Context.Res.Text("failed to send email: "+err.Error(), Context.Res.WithStatusCode(500))
		}

		return Context.Res.Json(Response{
			Success: true,
		})
	}

	return Context.Res.Text("Invalid path or method", Context.Res.WithStatusCode(404))
}
