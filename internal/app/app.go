package app

import (
	"context"
	"main/internal/env"
)

type App struct{}

func New() *App {
	_ = env.E //touch env package and env.init()
	return &App{}
}

func (a App) Run(ctx context.Context) {
	<-ctx.Done()
}
