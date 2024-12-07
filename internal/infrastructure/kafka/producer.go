package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
	keyID    uint
	ID       uint
}

func NewProducer(broker string, id uint) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"acks":              "all",
		"retries":           5,
	})
	if err != nil {
		return nil, err
	}
	return &Producer{
		producer: p,
		keyID:    0,
		ID:       id,
	}, nil
}

func (p *Producer) Produce(topic string, message []byte) error {
	p.keyID++
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
		Key:            []byte(fmt.Sprintf("%d-%d", p.keyID, p.ID)),
	}, nil)
}

func (p *Producer) Close() {
	p.producer.Flush(15 * 1000)
	p.producer.Close()
}
