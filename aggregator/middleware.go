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
			"func": "AggregateDistance",
		}).Info("AggregateDistance")
	}(time.Now())
	err = lm.next.AggregateDistance(d)
	return
}
