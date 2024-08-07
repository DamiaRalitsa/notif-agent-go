package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"gitlab.playcourt.id/notif-agent-go/helpers"
	"gitlab.playcourt.id/notif-agent-go/notification/mailer/model"
)

type SmtpClient interface {
	Notification
}

type gateway struct {
	BaseURL    string
	Host       string
	Port       string
	Username   string
	Password   string
	HttpClient *helpers.ToolsAPI
}

func (d gateway) SendEmail(ctx context.Context, to []string, subject string, message string, attachments []model.Attachments) (data interface{}, err error) {
	from := d.Username
	password := d.Password
	smtpHost := d.Host
	smtpPort := d.Port

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(to, ",")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = `multipart/mixed; boundary="MULTIPART_BOUNDARY"`

	header := ""
	for k, v := range headers {
		header += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	bodyHeader := "--MULTIPART_BOUNDARY\r\n" +
		`Content-Type: text/html; charset="UTF-8"` + "\r\n" +
		"Content-Transfer-Encoding: 7bit\r\n" +
		"\r\n" +
		message +
		"\r\n"

	newAttachments := ""
	for _, attachment := range attachments {
		newAttachments += "--MULTIPART_BOUNDARY\r\n" +
			`Content-Type: application/octet-stream` + "\r\n" +
			`Content-Transfer-Encoding: base64` + "\r\n" +
			`Content-Disposition: attachment; filename="` + attachment.FileName + `"` + "\r\n" +
			"\r\n" +
			base64.StdEncoding.EncodeToString(attachment.Content) +
			"\r\n"
	}

	newMessage := []byte(header + "\r\n" + bodyHeader + newAttachments + "--MULTIPART_BOUNDARY--")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, newMessage)

	if err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	fmt.Println("Email Sent Successfully!")
	return "OK", nil
}

func NewSmtpClient() SmtpClient {
	baseUrl := os.Getenv("FABD_API_CORE_URL")
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	httpClient := &helpers.ToolsAPI{}
	return &gateway{
		BaseURL:    baseUrl,
		Host:       host,
		Port:       port,
		HttpClient: httpClient,
		Username:   username,
		Password:   password,
	}
}
