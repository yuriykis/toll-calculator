package main

import (
	"fmt"

	"github.com/yuriykis/tolling/types"
)

type Aggregator interface {
	AggregateDistance(distance types.Distance) error
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
	return ia.store.Insert(distance)
}
