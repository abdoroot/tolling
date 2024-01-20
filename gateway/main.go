package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/abdoroot/tolling/aggregator/client"
	"github.com/sirupsen/logrus"
)

var (
	listenAddr  = ":6000"
	aggEndPoint = "http://127.0.0.1:3001/invoice"
)

type ApiFunc func(w http.ResponseWriter, r *http.Request) error

func MakeAPiFunc(fn ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			logrus.Error("MakeAPiFunc Eroor :", err)
			writeJson(w, http.StatusBadRequest, map[string]any{
				"err": err.Error(),
			})
		}
	}
}

func main() {
	aggClient := client.NewHTTP(aggEndPoint)
	inh := NewInvoiceHandler(aggClient)
	http.HandleFunc("/invoice", MakeAPiFunc(inh.handleGetInvoice))

	logrus.Infof("Gatawat running at port %v", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

type invoiceHandler struct {
	client *client.HTTPClient
}

func NewInvoiceHandler(c *client.HTTPClient) *invoiceHandler {
	return &invoiceHandler{
		client: c,
	}
}

func (h *invoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	values, ok := r.URL.Query()["obuid"]
	if !ok {
		return errors.New("please enter vlaid obuid")
	}
	obuIdString := values[0]
	obuIdInt, err := strconv.Atoi(obuIdString)
	if err != nil {
		return err
	}

	inv, err := h.client.GetInvoice(obuIdInt)
	if err != nil {
		return err
	}

	writeJson(w, http.StatusOK, inv)
	return nil
}

func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
