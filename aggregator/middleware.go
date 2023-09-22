package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"github.com/yuriykis/tolling/types"
)

type MetricsMiddleware struct {
	reqCounterAggregate prometheus.Counter
	reqCounterCalculate prometheus.Counter
	reqLatencyAggregate prometheus.Histogram
	reqLatencyCalculate prometheus.Histogram
	errCounterAggregate prometheus.Counter
	errCounterCalculate prometheus.Counter
	next                Aggregator
}

func NewMetricsMiddleware(next Aggregator) Aggregator {
	reqCounterAggregate := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})
	reqCounterCalculate := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
	})
	reqLatencyAggregate := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
	})
	reqLatencyCalculate := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculate",
	})
	errCounterAggregate := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregate",
	})
	errCounterCalculate := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate",
	})
	return &MetricsMiddleware{
		reqCounterAggregate: reqCounterAggregate,
		reqCounterCalculate: reqCounterCalculate,
		reqLatencyAggregate: reqLatencyAggregate,
		reqLatencyCalculate: reqLatencyCalculate,
		errCounterAggregate: errCounterAggregate,
		errCounterCalculate: errCounterCalculate,
		next:                next,
	}
}

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

func (mm *MetricsMiddleware) AggregateDistance(d types.Distance) (err error) {
	defer func(start time.Time) {
		mm.reqLatencyAggregate.Observe(time.Since(start).Seconds())
		mm.reqCounterAggregate.Inc()
		if err != nil {
			mm.errCounterAggregate.Inc()
		}
	}(time.Now())
	err = mm.next.AggregateDistance(d)
	return
}

func (mm *MetricsMiddleware) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		mm.reqLatencyCalculate.Observe(time.Since(start).Seconds())
		mm.reqCounterCalculate.Inc()
		if err != nil {
			mm.errCounterCalculate.Inc()
		}
	}(time.Now())
	invoice, err = mm.next.CalculateInvoice(obuID)
	return
}
