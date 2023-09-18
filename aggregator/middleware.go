package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yuriykis/tolling/types"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (lm *LogMiddleware) AggregateDistance(d types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())
	err = lm.next.AggregateDistance(d)
	return
}

func (lm *LogMiddleware) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			dist   float64
			amount float64
		)
		if invoice != nil {
			dist = invoice.TotalDistance
			amount = invoice.TotoalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took":   time.Since(start),
			"err":    err,
			"obuID":  obuID,
			"dist":   dist,
			"amount": amount,
		}).Info("CalculateInvoice")
	}(time.Now())
	invoice, err = lm.next.CalculateInvoice(obuID)
	return
}
