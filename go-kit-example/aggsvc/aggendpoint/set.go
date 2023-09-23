package aggendpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/yuriykis/tolling/go-kit-example/aggsvc/aggservice"
	"github.com/yuriykis/tolling/types"
)

type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

type AggregateRequest struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuID"`
	Unix  int64   `json:"unix"`
}

type AggregateResponse struct {
	Error error `json:"error,omitempty"`
}

type CalculateRequest struct {
	OBUID int `json:"obuID"`
}

type CalculateResponse struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotoalAmount  float64 `json:"totalAmount"`
	Error         error   `json:"error,omitempty"`
}

func makeAggregateEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AggregateRequest)
		err := s.Aggregate(ctx, types.Distance{
			Value: req.Value,
			OBUID: req.OBUID,
			Unix:  req.Unix,
		})
		return AggregateResponse{Error: err}, nil
	}
}

func makeCalculateEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CalculateRequest)
		inv, err := s.Calculate(ctx, req.OBUID)
		if err != nil {
			return CalculateResponse{Error: err}, nil
		}
		return CalculateResponse{
			OBUID:         inv.OBUID,
			TotalDistance: inv.TotalDistance,
			TotoalAmount:  inv.TotoalAmount,
		}, nil
	}
}

func (s Set) Aggregate(ctx context.Context, dist types.Distance) error {
	_, err := s.AggregateEndpoint(ctx, AggregateRequest{
		Value: dist.Value,
		OBUID: dist.OBUID,
		Unix:  dist.Unix,
	})
	return err
}

func (s Set) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	resp, err := s.CalculateEndpoint(ctx, CalculateRequest{
		OBUID: obuID,
	})
	if err != nil {
		return nil, err
	}
	response := resp.(CalculateResponse)
	return &types.Invoice{
		OBUID:         response.OBUID,
		TotalDistance: response.TotalDistance,
		TotoalAmount:  response.TotoalAmount,
	}, nil

}
