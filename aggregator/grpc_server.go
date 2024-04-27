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
	return nil, g.svc.AggregateDistance(dist)
}

func (g *GRPCDistanceAggregatorServer) GetInvoice(ctx context.Context, req *types.GetInvoiceRequest) (*types.GetInvoiceResponse, error) {
	inv, err := g.svc.Invoice(int(req.ObuID))
	if err != nil {
		return nil, err
	}
	return &types.GetInvoiceResponse{
		ObuID:         int64(inv.OBUID),
		InvoiceAmount: inv.InvoiceAmount,
		TotalDistance: inv.TotalDistance,
	}, nil
}
