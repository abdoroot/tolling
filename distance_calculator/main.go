package main

import (
	"fmt"

	"github.com/abdoroot/tolling/aggregator/client"
)

const (
	topic              = "data"
	aggregatorEndPoint = "http://127.0.0.1:3001/aggregate"
)

func main() {
	dc := client.New(aggregatorEndPoint)
	svc := NewCalculateService()
	svc = NewLogMiddleware(svc)
	c, err := NewKafkaConsumer(topic, svc, dc)
	if err != nil {
		fmt.Println(err)
	}
	c.Run()
}
