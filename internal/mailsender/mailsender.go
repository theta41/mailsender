package mailsender

import (
	"context"
	"main/internal/model"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type MailSender struct {
	workers int
	msgPool chan model.Email
	stop    chan struct{}
	wg      sync.WaitGroup
}

func New(workers int) *MailSender {
	m := &MailSender{
		workers: workers,
		msgPool: make(chan model.Email, workers),
		stop:    make(chan struct{}),
	}
	m.run()
	return m
}

func (m *MailSender) run() {
	for i := 1; i < m.workers; i++ {
		m.wg.Add(1)
		go m.worker()
	}
}

func (m *MailSender) Stop() {
	close(m.stop)
	m.wg.Wait()
	close(m.msgPool)
}

func (m *MailSender) worker() {
	defer m.wg.Done()
	for {
		select {
		case msg := <-m.msgPool:
			printEmail(msg)
		case <-m.stop:
			return
		}
	}
}

func printEmail(msg model.Email) {
	time.Sleep(250 * time.Millisecond)
	logrus.Infof("printMessage %+v", msg)
}

func (m *MailSender) SendEmail(ctx context.Context, msg model.Email) error {
	select {
	case m.msgPool <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
