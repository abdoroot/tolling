package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/abdoroot/tolling/internal/envutil"
	"github.com/abdoroot/tolling/types"
	"github.com/gorilla/websocket"
)

var sendInterval = envutil.Duration("OBU_SEND_INTERVAL", 5*time.Second)

var receiverEndPoint = envutil.String("OBU_WS_ENDPOINT", "ws://localhost:3000")

func genCord() float64 {
	nf := float64(rand.Intn(100) + 1)
	return rand.Float64() + nf
}

func genLatlog() (float64, float64) {
	return genCord(), genCord()
}

func main() {
	obuids := genOBUIDs(envutil.Int("OBU_COUNT", 20))
	conn, _, err := websocket.DefaultDialer.Dial(receiverEndPoint, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()
	for {
		for i := 0; i < len(obuids); i++ {
			lat, long := genLatlog()
			data := types.OBUdata{
				OBUID: obuids[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal("WriteJSON:", err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func genOBUIDs(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}
