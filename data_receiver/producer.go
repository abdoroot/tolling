package main

import (
	"encoding/json"
	"fmt"

	"github.com/abdoroot/tolling/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type DataProducer interface {
	Produce() error
}

type kafkaProducer struct {
	Producer *kafka.Producer
}

func NewkafkaProducer() (*kafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &kafkaProducer{
		Producer: p,
	}, nil
}

func (p *kafkaProducer) Produce(data types.OBUdata, topic string) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: b,
	}, nil)
}
