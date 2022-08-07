package main

import (
	"context"
	"main/internal/app"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	ServiceName    = "mailsender"
	ServiceVersion = "v1.0.0"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	a := app.New() //init app and load env

	logrus.Infof("Starting %v %v", ServiceName, ServiceVersion)

	a.Run(ctx)

	logrus.Infof("Stopped %v %v", ServiceName, ServiceVersion)
}
