package main

import (
	"fmt"

	"github.com/abdoroot/tolling/aggregator/client"
	"github.com/abdoroot/tolling/internal/envutil"
)

const (
	defaultTopic              = "data"
	defaultKafkaBootstrap     = "localhost:9092"
	defaultKafkaGroupID       = "myGroup"
	defaultAggregatorGRPCAddr = "127.0.0.1:3002"
)

func main() {
	topic := envutil.String("KAFKA_TOPIC", defaultTopic)
	kafkaBootstrapServers := envutil.String("KAFKA_BOOTSTRAP_SERVERS", defaultKafkaBootstrap)
	kafkaGroupID := envutil.String("KAFKA_GROUP_ID", defaultKafkaGroupID)
	aggGRPCEndPoint := envutil.String("AGGREGATOR_GRPC_ENDPOINT", defaultAggregatorGRPCAddr)

	GRPCClient, err := client.NewGRPC(aggGRPCEndPoint)
	if err != nil {
		fmt.Println(err)
	}
	svc := NewCalculateService()
	svc = NewLogMiddleware(svc)
	c, err := NewKafkaConsumer(kafkaBootstrapServers, kafkaGroupID, topic, svc, GRPCClient)
	if err != nil {
		fmt.Println(err)
	}
	c.Run()
}
