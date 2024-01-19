package main

import (
	"context"

	"github.com/abdoroot/tolling/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcAggregaorSever struct {
	types.AggreagatorServer
	srv Aggregator
}

func NewGrpcAggregaorSever(srv Aggregator) *GrpcAggregaorSever {
	return &GrpcAggregaorSever{
		srv: srv,
	}
}

func (g *GrpcAggregaorSever) AggregateDistance(ctx context.Context, req *types.DistanceRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(req.OBUID),
		Value: req.Value,
		Unix:  req.Unix,
	}

	if err := g.srv.AggregateDistance(distance); err != nil {
		return &types.None{}, status.Error(codes.Internal, "grpc aggreagate error"+err.Error())
	}
	return &types.None{}, nil
}
