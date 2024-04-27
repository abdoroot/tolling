package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/abdoroot/tolling/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	HTTPlistenAddr   = ":3001"
	GRPClistenAddr   = ":3002"
	MatrixlistenAddr = ":2112"
	store            = NewMemoryStore()
	srv              = NewInvoiceAggregator(store) //service
)

func main() {
	srv = NewLogMiddleware(srv)
	srv = NewMetricsMiddleWare(srv)
	//matrix http server
	go makeMartixServer()
	//Http transport
	go makeHttpTransport(srv, HTTPlistenAddr)
	makeGRPCTransport(srv, GRPClistenAddr)
}

func makeGRPCTransport(srv Aggregator, listenAddr string) error {
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	fmt.Println("grpc server running at port", listenAddr)
	grpcServer := grpc.NewServer()
	types.RegisterAggreagatorServer(grpcServer, NewGrpcAggregaorSever(srv))
	grpcServer.Serve(l)
	return nil
}

func makeMartixServer() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	fmt.Println("martex running at port", MatrixlistenAddr)
	http.ListenAndServe(MatrixlistenAddr, mux)
}

func makeHttpTransport(srv Aggregator, listenAddr string) {
	fmt.Println("http transport running at port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(srv))
	http.HandleFunc("/invoice", handleGetInvoice(srv))
	http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(srv Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obuid"]
		if !ok {
			writeJson(w, http.StatusBadRequest, map[string]any{
				"error": "please enter vlaid obuid",
			})
			return
		}

		obuid, err := strconv.Atoi(values[0])
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]any{
				"error": "please enter vlaid obuid",
			})
			return
		}

		inv, err := srv.CalculateInvoice(obuid)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]any{
				"error": err.Error(),
			})
			return
		}

		writeJson(w, http.StatusOK, inv)
	}
}

func handleAggregate(srv Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := types.Distance{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			writeJson(w, http.StatusBadRequest, map[string]any{
				"error": "json decode error" + err.Error(),
			})
			return
		}
		if err := srv.AggregateDistance(data); err != nil {
			return
		}
	}
}

func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
