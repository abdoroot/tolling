package main

import (
	"fmt"

	"github.com/abdoroot/tolling/aggregator/client"
)

const (
	topic           = "data"
	aggHTTPEndPoint = "http://127.0.0.1:3001/aggregate"
	aggGRPCEndPoint = "127.0.0.1:3002"
)

func main() {
	//HTTPClient := client.NewHTTP(aggHTTPEndPoint)
	GRPCClient, err := client.NewGRPC(aggGRPCEndPoint)
	if err != nil {
		fmt.Println(err)
	}
	svc := NewCalculateService()
	svc = NewLogMiddleware(svc)
	c, err := NewKafkaConsumer(topic, svc, GRPCClient)
	if err != nil {
		fmt.Println(err)
	}
	c.Run()
}
