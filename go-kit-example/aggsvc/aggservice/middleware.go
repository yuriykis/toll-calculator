package aggservice

import (
	"context"

	"github.com/yuriykis/tolling/types"
)

type Middleware func(Service) Service

type loginMiddleware struct {
	next Service
}

func newLoginMiddleware() Middleware {
	return func(next Service) Service {
		return &loginMiddleware{
			next: next,
		}
	}
}

func (mw *loginMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	// TODO implement loginMiddleware.Aggregate
	return nil
}

func (mw *loginMiddleware) Calculate(
	_ context.Context,
	obuID int,
) (*types.Invoice, error) {
	// TODO implement loginMiddleware.Calculate
	return nil, nil
}

type instrumentationMiddleware struct {
	next Service
}

func newInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return &instrumentationMiddleware{
			next: next,
		}
	}
}

func (mw *instrumentationMiddleware) Aggregate(
	ctx context.Context,
	dist types.Distance,
) error {
	// TODO implement instrumentationMiddleware.Aggregate
	return nil
}

func (mw *instrumentationMiddleware) Calculate(
	ctx context.Context,
	obuID int,
) (*types.Invoice, error) {
	// TODO implement instrumentationMiddleware.Calculate
	return nil, nil
}
