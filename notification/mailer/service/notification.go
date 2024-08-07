package service

import (
	"context"

	"gitlab.playcourt.id/notif-agent-go/notification/mailer/model"
)

type Notification interface {
	Send(ctx context.Context, body model.Mail) (data interface{}, err error)
}
