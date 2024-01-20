package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/abdoroot/tolling/types"
	"github.com/sirupsen/logrus"
)

type AggregatorClient interface {
	AggregateInvoice(types.Distance) error
	GetInvoice(id int) (*types.Invoice, error)
}

type HTTPClient struct {
	endPoint string
}

func NewHTTP(endPoint string) *HTTPClient {
	return &HTTPClient{
		endPoint: endPoint,
	}
}

func (c *HTTPClient) GetInvoice(ObuId int) (*types.Invoice, error) {
	ObuIdString := strconv.Itoa(ObuId)
	reqUrl := strings.Join([]string{
		c.endPoint,
		"?obuid=",
		ObuIdString,
	}, "")
	logrus.Infof("Get Invoice Client Request Url :%v", reqUrl)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("the service response with non 200 status code")
	}

	inv := &types.Invoice{}
	if err := json.NewDecoder(resp.Body).Decode(inv); err != nil {
		return nil, err
	}
	return inv, nil
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
