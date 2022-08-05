package looper

import (
	"context"
	"main/internal/model"
)

type MessageReader interface {
	ReadMessage(context.Context) (model.Message, error)
}
