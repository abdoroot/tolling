package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/abdoroot/tolling/types"
)

type AggregatorClient interface {
	AggregateInvoice(types.Distance) error
}

type HTTPClient struct {
	endPoint string
}

func NewHTTP(endPoint string) *HTTPClient {
	return &HTTPClient{
		endPoint: endPoint,
	}
}

func (c *HTTPClient) AggregateInvoice(data types.Distance) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.endPoint, bytes.NewReader(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("the service response with non 200 status code")
	}

	return nil
}
