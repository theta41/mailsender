package looper

import (
	"context"
	"errors"
	"io"

	"github.com/sirupsen/logrus"
)

type readerSender struct {
	reader MessageReader
	sender EmailSender
}

type status int

const (
	statusSent status = iota
	statusRetry
	statusCanceled
)

func (s *readerSender) readAndSend(ctx context.Context, log *logrus.Entry) status {

	m, err := s.reader.ReadMessage(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return statusCanceled
		}
		if errors.Is(err, io.EOF) {
			log.Info("reader close detected")
			return statusCanceled
		}
		log.WithError(err).Error("cannot read message")
		return statusRetry
	}

	email := m.ToEmail()
	err = s.sender.SendEmail(ctx, email)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Errorf("sending canceled, email lost %+v", email)
			return statusCanceled
		}
		log.WithError(err).Errorf("cannot send, email lost %+v", email)
		return statusRetry
	}

	log.WithFields(logrus.Fields{
		"email":   email.Address,
		"content": email.Content,
	}).Info("successfully sent message")

	return statusSent
}
