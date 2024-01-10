package main

import "fmt"

const topic = "data"

func main() {
	svc := NewCalculateService()
	svc = NewLogMiddleware(svc)
	c, err := NewKafkaConsumer(topic, svc)
	if err != nil {
		fmt.Println(err)
	}
	c.Run()
}
