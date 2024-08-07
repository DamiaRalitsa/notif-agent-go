package service

import (
	"context"

	"gitlab.playcourt.id/notif-agent-go/notification/mailer/model"
)

type Notification interface {
	SendEmail(ctx context.Context, to []string, subject string, message string, attachments []model.Attachments) (data interface{}, err error)
}
