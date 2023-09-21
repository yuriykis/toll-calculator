package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/yuriykis/tolling/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) Aggregate(
	ctx context.Context,
	aggReq *types.AggregateRequest,
) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		c.Endpoint+"/aggregate",
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	defer resp.Body.Close()

	return nil
}

func (c *HTTPClient) GetInvoice(ctx context.Context, obuID int) (*types.Invoice, error) {
	invReq := types.GetInvoiceRequest{
		ObuId: int64(obuID),
	}
	b, err := json.Marshal(&invReq)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/%s?obu=%d", c.Endpoint, "invoice", obuID)
	logrus.Infof("requesting invoice from %s", endpoint)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	var inv types.Invoice
	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &inv, nil
}
