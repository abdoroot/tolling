package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/abdoroot/tolling/types"
)

func main() {
	var (
		listenAddr = ":3001"
		store      = NewMemoryStore()
		srv        = NewInvoiceAggregator(store) //service
	)
	srv = NewLogMiddleware(srv)
	//Http transport
	makeHttpTransport(srv, listenAddr)
}

func makeHttpTransport(srv Aggregator, listenAddr string) {
	fmt.Println("http transport running at port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(srv))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(srv Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := types.Distance{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("content-type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"error": "json decode error" + err.Error(),
			})
			return
		}
		if err := srv.AggregateDistance(data); err != nil {
			return
		}
	}
}
