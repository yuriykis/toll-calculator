package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/yuriykis/tolling/types"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(distance types.Distance) error
	CalculateInvoice(obuID int) (*types.Invoice, error)
}
type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (ia *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Printf("Aggregating and inserting distance %v\n", distance)
	logrus.WithFields(logrus.Fields{
		"obu_id": distance.OBUID,
		"value":  distance.Value,
		"unix":   distance.Unix,
	}).Info("Aggregating and inserting distance")
	return ia.store.Insert(distance)
}

func (ia *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {
	dist, err := ia.store.Get(obuID)
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
