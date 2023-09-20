package main

import (
	"context"

	"github.com/yuriykis/tolling/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregatorServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

func (s *GRPCAggregatorServer) Aggregate(
	ctx context.Context,
	req *types.AggregateRequest,
) (*types.None, error) {
	dist := types.Distance{
		OBUID: int(req.ObuId),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return &types.None{}, s.svc.AggregateDistance(dist)
}

func (s *GRPCAggregatorServer) GetInvoice(
	ctx context.Context,
	req *types.GetInvoiceRequest,
) (*types.Invoice, error) {
	// not implemented
	return nil, nil
}
