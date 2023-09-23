package aggservice

import (
	"context"

	"github.com/go-kit/log"
	"github.com/yuriykis/tolling/types"
)

const basePrice = 3.15

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

type BasicService struct {
	store Storer
}

func newBasicService(store Storer) Service {
	return &BasicService{
		store: store,
	}
}

func (svc *BasicService) Aggregate(_ context.Context, dist types.Distance) error {
	return svc.store.Insert(dist)
}

func (svc *BasicService) Calculate(_ context.Context, obuID int) (*types.Invoice, error) {
	dist, err := svc.store.Get(obuID)
	if err != nil {
		return nil, err
	}
	invoice := &types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotoalAmount:  basePrice + dist,
	}
	return invoice, nil
}

// NewAggregatorService will construct a complete microservice with logging and instrumentation middleware
func New(logger log.Logger) Service {
	var svc Service
	{

		svc = newBasicService(NewMemoryStore())
		svc = newLoggingMiddleware(logger)(svc)
		svc = newInstrumentationMiddleware()(svc)

	}
	return svc
}
