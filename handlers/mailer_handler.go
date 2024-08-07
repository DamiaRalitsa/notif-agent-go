package handlers

import (
	"bytes"
	"context"
	"io/ioutil"
	"mime"
	"path/filepath"
	"text/template"

	"gitlab.playcourt.id/notif-agent-go/notification/mailer/model"
	"gitlab.playcourt.id/notif-agent-go/notification/mailer/service"
)

type MailerHandler struct {
	notificationService service.Notification
}

func NewMailerHandler(notificationService service.Notification) *MailerHandler {
	return &MailerHandler{
		notificationService: notificationService,
	}
}

func (h *MailerHandler) SendEmailWithAttachments(ctx context.Context, to []string, subject string, emailTemplate string, templateData interface{}, filePaths []string) (data interface{}, err error) {
	// Read the HTML template file
	htmlTemplate, err := ioutil.ReadFile(emailTemplate)
	if err != nil {
		return nil, err
	}

	// Parse the HTML template
	tmpl, err := template.New("emailTemplate").Parse(string(htmlTemplate))
	if err != nil {
		return nil, err
	}

	// Execute the template with the provided data
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, templateData); err != nil {
		return nil, err
	}

	// Set the resulting string as the message content
	message := tpl.String()

	attachments := make([]model.Attachments, 0)

	for _, filePath := range filePaths {
		// Read the file content
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		// Determine the file's content type
		fileName := filepath.Base(filePath)
		contentType := mime.TypeByExtension(filepath.Ext(fileName))

		// Map the file content and metadata into the model.Attachments struct
		attachment := model.Attachments{
			FileName:    fileName,
			Content:     fileContent,
			Encoding:    "base64", // Assuming base64 encoding
			ContentType: contentType,
		}

		attachments = append(attachments, attachment)
	}

	// Create the Mail struct
	mail := model.Mail{
		To:          to,
		Subject:     subject,
		Message:     message,
		Attachments: attachments,
	}

	// Call the existing SendEmail method with the Mail struct
	return h.notificationService.SendEmail(ctx, mail)
}
