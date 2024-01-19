package client

import (
	"context"

	"github.com/abdoroot/tolling/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	endPoint string
	types.AggreagatorClient
}

func NewGRPC(endPoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endPoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := types.NewAggreagatorClient(conn)
	return &GRPCClient{
		endPoint:          endPoint,
		AggreagatorClient: c,
	}, nil
}

func (c *GRPCClient) AggregateInvoice(data types.Distance) error {
	req := types.DistanceRequest{
		OBUID: int64(data.OBUID),
		Value: data.Value,
		Unix:  data.Unix,
	}
	resp, err := c.AggregateDistance(context.TODO(), &req)
	_ = resp
	return err
}
