package main

import (
	"encoding/json"

	"github.com/abdoroot/tolling/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type DataProducer interface {
	ProduceData(types.OBUdata) error
}

type kafkaProducer struct {
	Producer *kafka.Producer
	Topic    string
}

func NewkafkaProducer(topic string) (*kafkaProducer, error) {
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
					//fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					//fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &kafkaProducer{
		Producer: p,
		Topic:    topic,
	}, nil
}

func (p *kafkaProducer) ProduceData(data types.OBUdata) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: b,
	}, nil)
}
