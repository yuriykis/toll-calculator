package aggendpoint

import (
	"context"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	log "github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"github.com/yuriykis/tolling/go-kit-example/aggsvc/aggservice"
	"github.com/yuriykis/tolling/types"
	"golang.org/x/time/rate"
)

type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

func New(svc aggservice.Service, logger log.Logger) Set {
	var aggregateEndpoint endpoint.Endpoint
	{
		aggregateEndpoint = makeAggregateEndpoint(svc)
		// Sum is limited to 1 request per second with burst of 1 request.
		// Note, rate is defined as a time interval between requests.
		aggregateEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			aggregateEndpoint,
		)
		aggregateEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			aggregateEndpoint,
		)

		// aggregateEndpoint = LoggingMiddleware(
		// 	log.With(logger, "method", "Sum"),
		// )(
		// 	aggregateEndpoint,
		// )
		// aggregateEndpoint = InstrumentingMiddleware(
		// 	duration.With("method", "Sum"),
		// )(
		// 	aggregateEndpoint,
		// )
	}
	var calculateEndpoint endpoint.Endpoint
	{
		calculateEndpoint = makeCalculateEndpoint(svc)
		// Concat is limited to 1 request per second with burst of 100 requests.
		// Note, rate is defined as a number of requests per second.
		calculateEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 100),
		)(
			calculateEndpoint,
		)
		calculateEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			calculateEndpoint,
		)
		// calculateEndpoint = LoggingMiddleware(
		// 	log.With(logger, "method", "Concat"),
		// )(
		// 	calculateEndpoint,
		// )
		// calculateEndpoint = InstrumentingMiddleware(
		// 	duration.With("method", "Concat"),
		// )(
		// 	calculateEndpoint,
		// )
	}
	return Set{
		AggregateEndpoint: aggregateEndpoint,
		CalculateEndpoint: calculateEndpoint,
	}
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
