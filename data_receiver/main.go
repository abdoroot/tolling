package main

import (
	"log"
	"net/http"

	"github.com/abdoroot/tolling/types"
	"github.com/gorilla/websocket"
)

var (
	Topic    = "data"
	upgrader = websocket.Upgrader{}
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
		dr.producer.ProduceData(data)
	}
}

func main() {
	var p DataProducer
	p, err := NewkafkaProducer(Topic)
	if err != nil {
		log.Fatal(err)
	}

	//producer middleware
	pm := NewLoggingMiddleware(p)
	dr := NewDataReceiver(pm)

	http.HandleFunc("/", dr.handleWc)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
