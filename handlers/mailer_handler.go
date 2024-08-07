package handlers

import (
	"context"
	"io/ioutil"
	"mime"
	"path/filepath"

	"gitlab.playcourt.id/notif-agent-go/notification/mailer/model"
	"gitlab.playcourt.id/notif-agent-go/notification/mailer/service"
)

type MailerHandler struct {
	notificationService service.Notification
}

func (h *MailerHandler) SendEmailWithAttachments(ctx context.Context, to []string, subject string, message string, filePaths []string) (data interface{}, err error) {
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

	// Call the existing SendEmail method with the attachments
	return h.notificationService.SendEmail(ctx, to, subject, message, attachments)
}
