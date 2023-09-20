package client

import (
	"context"

	"github.com/yuriykis/tolling/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
