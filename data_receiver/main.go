package main

import (
	"log"
	"net/http"

	"github.com/abdoroot/tolling/internal/envutil"
	"github.com/abdoroot/tolling/types"
	"github.com/gorilla/websocket"
)

var (
	topic                 = envutil.String("KAFKA_TOPIC", "data")
	kafkaBootstrapServers = envutil.String("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
	listenAddr            = envutil.String("DATA_RECEIVER_LISTEN_ADDR", ":3000")
	upgrader              = websocket.Upgrader{}
)

type DataReceiver struct {
	msgch    chan *types.OBUdata
	producer DataProducer
}

func NewDataReceiver(p DataProducer) *DataReceiver {
	return &DataReceiver{
		msgch:    make(chan *types.OBUdata, 128),
		producer: p,
	}
}

func (dr *DataReceiver) handleWc(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	go dr.reciveLoop(c)
}

func (dr *DataReceiver) reciveLoop(c *websocket.Conn) {
	for {
		data := types.OBUdata{}
		if err := c.ReadJSON(&data); err != nil {
			log.Println("read:", err)
			c.Close()
			continue
		}
		err := dr.producer.ProduceData(data)
		if err != nil {
			log.Println("produce data error", err)
			continue
		}
	}
}

func main() {
	var p DataProducer
	p, err := NewkafkaProducer(kafkaBootstrapServers, topic)
	if err != nil {
		log.Fatal(err)
	}

	//producer middleware
	pm := NewLoggingMiddleware(p)
	dr := NewDataReceiver(pm)

	http.HandleFunc("/", dr.handleWc)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
