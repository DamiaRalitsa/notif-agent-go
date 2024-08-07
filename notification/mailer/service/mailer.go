package service

import (
	"context"

	"gitlab.playcourt.id/notif-agent-go/notification/mailer/model"
)

type Notification interface {
	SendEmail(ctx context.Context, mail model.Mail) (data interface{}, err error)
}
