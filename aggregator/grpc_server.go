package main

import (
	"context"

	"github.com/swarajroy/toll_calculator/types"
)

type GRPCDistanceAggregatorServer struct {
	types.UnimplementedDistanceAggregatorServer
	svc Aggregator
}

func NewGRPCDistanceAggregatorServer(svc Aggregator) *GRPCDistanceAggregatorServer {
	return &GRPCDistanceAggregatorServer{
		svc: svc,
	}
}

func (g *GRPCDistanceAggregatorServer) AggregateDistance(ctx context.Context, req *types.AggregatorDistanceRequest) (*types.None, error) {
	dist := types.Distance{
		OBUID: int(req.ObuID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	g.svc.AggregateDistance(dist)
	return nil, nil
}
