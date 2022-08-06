package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"main/internal/model"

	"github.com/sirupsen/logrus"

	"github.com/segmentio/kafka-go"
)

const (
	keyEmail = "email"
)

type mqEmail struct {
	Address string `json:"address"`
	Content string `json:"content"`
}

type Client struct {
	reader *kafka.Reader
}

func NewClient(brokers []string, topic, groupId string) (*Client, error) {
	if len(brokers) == 0 || brokers[0] == "" {
		return nil, errors.New("invalid kafka configuration")
	}

	c := Client{}

	log := logrus.WithField("mq", "kafka-reader")

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupId,
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e6,
		//Logger:      kafka.LoggerFunc(log.Infof),
		ErrorLogger: kafka.LoggerFunc(log.Errorf),
	})

	return &c, nil
}

func (c *Client) Close() {
	logrus.Info("close kafka reader")
	c.reader.Close()
}

func (c *Client) ReadMessage(ctx context.Context) (model.Message, error) {
	log := logrus.WithField("mq", "kafka-ReadMessage")
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return model.Message{}, err
		}

		key := string(m.Key)
		log := log.WithField("key", key)

		switch key {
		default:
			log.Error("unknown message key")

		case keyEmail:
			var msg mqEmail
			err := json.Unmarshal(m.Value, &msg)
			if err != nil {
				log.Error("unmarshal value error: ", err)
				continue
			}

			log.Info("got message: ", msg)

			return model.Message(msg), nil
		}
	}
}
