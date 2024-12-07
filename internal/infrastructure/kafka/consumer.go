package kafka

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
}

func NewConsumer(broker string, group string) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          group,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
	}, nil
}

func (c *Consumer) SubscribeTopics(topics []string) error {
	return c.consumer.SubscribeTopics(topics, nil)
}

func (c *Consumer) Consume(timeout time.Duration) (*kafka.Message, error) {
	return c.consumer.ReadMessage(timeout)
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func (c *Consumer) ConsumeMessagesEternally(dataCallback func([]byte) error, successCallback func(string, ...interface{}), errorCallback func(string, ...interface{})) {
	if dataCallback == nil {
		return
	}

	safeCallback := func(callback func(string, ...interface{}), format string, v ...interface{}) {
		if callback != nil {
			callback(format, v...)
		}
	}

	// TODO: redis
	messageKeys := make(map[string]bool)

	for {
		msg, err := c.Consume(time.Second)
		if err == nil {
			keyStr := string(msg.Key)
			if _, ok := messageKeys[keyStr]; ok {
				continue
			}
			messageKeys[keyStr] = true

			err := dataCallback(msg.Value)

			if err != nil {
				safeCallback(errorCallback, "got an error while consuming messages: %s", err.Error())
			} else {
				safeCallback(successCallback, "message consumed successfully")
			}
		} else if !err.(kafka.Error).IsTimeout() {
			go safeCallback(errorCallback, "consumer error: %s", err.Error())
		}
	}
}
