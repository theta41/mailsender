package looper

import (
	"context"
	"main/internal/model"
)

type EmailSender interface {
	SendEmail(context.Context, model.Email) error
}
