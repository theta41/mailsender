package app

import (
	"context"
	"main/internal/env"
	"main/internal/kafka"
	"main/internal/looper"
	"main/internal/mailer"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct{}

func New() App {
	_ = env.E //touch env package and env.init()
	return App{}
}

func (a App) Run(ctx context.Context) {

	log := logrus.WithField("stage", "app-Run")

	reader, err := kafka.NewClient(
		env.E.C.Kafka.Brokers,
		env.E.C.Kafka.Topic,
		env.E.C.Kafka.GroupId,
	)
	if err != nil {
		log.WithError(err).Panic("cannot create kafka client")
		return
	}
	defer reader.Close()

	sender := mailer.New(env.E.C.RequstsInPeriod)
	defer sender.Stop()

	looper := looper.New(
		time.Duration(env.E.C.PeriodMinutes)*time.Minute,
		env.E.C.RequstsInPeriod,
		reader,
		sender,
	)

	looper.Loop(ctx)
}
