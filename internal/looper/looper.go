package looper

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type Looper interface {
	Loop(ctx context.Context)
}

type looperLimited struct {
	period   time.Duration
	requests int
	readerSender
}

type looper struct {
	readerSender
}

func New(period time.Duration, requests int, reader MessageReader, sender EmailSender) Looper {

	if period > 0 && requests > 0 {
		return &looperLimited{
			period:       period,
			requests:     requests,
			readerSender: readerSender{reader, sender},
		}

	}

	return &looper{
		readerSender: readerSender{reader, sender},
	}
}

func (o *looper) Loop(ctx context.Context) {
	log := logrus.WithField("stage", "looper")

	for o.readAndSend(ctx, log) != statusCanceled {
	}
}

func (o *looperLimited) Loop(ctx context.Context) {
	ticker := time.NewTicker(o.period)
	log := logrus.WithField("stage", "looperLimited")

	for {
		select {
		case <-ticker.C:

			counter := 0
			for counter < o.requests {

				result := o.readAndSend(ctx, log)
				if result == statusCanceled {
					return
				}
				if result == statusSent {
					counter++
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
