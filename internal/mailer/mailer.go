package mailer

import (
	"context"
	"main/internal/model"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Mailer struct {
	workers int
	msgPool chan model.Email
	stop    chan struct{}
	wg      sync.WaitGroup
}

func New(workers int) *Mailer {
	m := &Mailer{
		workers: workers,
		msgPool: make(chan model.Email, workers),
		stop:    make(chan struct{}),
	}
	m.start()
	return m
}

func (m *Mailer) start() {
	logrus.Infof("start mailer workers")
	for i := 1; i < m.workers; i++ {
		m.wg.Add(1)
		go m.worker()
	}
}

func (m *Mailer) Stop() {
	logrus.Info("stopping mailer workers")
	close(m.stop)
	m.wg.Wait()
	close(m.msgPool)
	logrus.Info("all mailer workers are stopped")
}

func (m *Mailer) worker() {
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

func (m *Mailer) SendEmail(ctx context.Context, msg model.Email) error {
	select {
	case m.msgPool <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
