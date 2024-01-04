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
	msgch chan *types.OBUdata
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan *types.OBUdata, 128),
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
	dataProducer, err := NewkafkaProducer()
	if err != nil {
		log.Fatal(err)
	}
	for {
		data := types.OBUdata{}
		if err := c.ReadJSON(&data); err != nil {
			log.Println("read:", err)
			c.Close()
			continue
		}
		dataProducer.Produce(data, Topic)
	}
}

func main() {
	dr := NewDataReceiver()
	http.HandleFunc("/", dr.handleWc)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
