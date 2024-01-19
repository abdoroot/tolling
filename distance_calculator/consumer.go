package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/abdoroot/tolling/aggregator/client"
	"github.com/abdoroot/tolling/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculateServicer
	aggClient   *client.Client
}

func NewKafkaConsumer(topic string, svc CalculateServicer, aggClient *client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	if err := c.SubscribeTopics([]string{topic}, nil); err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
		aggClient:   aggClient,
	}, nil
}

func (k *KafkaConsumer) Run() {
	fmt.Println("KafkaConsumer is running")
	k.isRunning = true
	k.ReadMessageLoop()
}

func (k *KafkaConsumer) ReadMessageLoop() {
	for k.isRunning {
		msg, err := k.consumer.ReadMessage(-1)
		if err != nil {
			fmt.Println("kafka consume err ", err)
			continue
		}
		data := types.OBUdata{}
		if err := json.NewDecoder(bytes.NewReader(msg.Value)).Decode(&data); err != nil {
			fmt.Println("Decode json error :", err)
			continue
		}

		distance, err := k.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Error("CalculateDistance err", err)
			continue
		}

		req := types.Distance{
			OBUID: data.OBUID,
			Value: distance,
			Unix:  time.Now().Unix(),
		}

		if err := k.aggClient.AggregateInvoice(req); err != nil {
			logrus.Error("aggregate invoice err", err)
			continue
		}
	}
}
