package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	syncProducer sarama.SyncProducer
	topic        string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		syncProducer: producer,
		topic:        topic,
	}, nil
}

func (p *Producer) SendMessage(key string, value interface{}) error {
	// JSON'a çevir
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Kafka mesajı hazırla
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(jsonValue),
	}

	// Kafka'ya gönder
	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message sent to partition %d offset %d", partition, offset)
	return nil
}

func (p *Producer) Close() error {
	return p.syncProducer.Close()
}
