package client

import (
	"context"

	"github.com/yuriykis/tolling/types"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	Endpoint string
	client   types.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	// conn, err := grpc.Dial(
	// 	endpoint,
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )
	conn, err := grpc.Dial(
		endpoint,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (c GRPCClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	_, err := c.client.Aggregate(ctx, aggReq)
	return err
}

func (c GRPCClient) GetInvoice(ctx context.Context, obuID int) (*types.Invoice, error) {
	// not implemented
	return nil, nil
}
