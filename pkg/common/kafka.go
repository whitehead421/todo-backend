package common

import (
	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter(env *Environment) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{env.KafkaBrokers},
		Topic:    env.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})
}

func NewKafkaReader(env *Environment) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{env.KafkaBrokers},
		Topic:    env.KafkaTopic,
		GroupID:  env.KafkaGroupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}
