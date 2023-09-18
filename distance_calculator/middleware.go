package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yuriykis/tolling/types"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LogMiddleware{
		next: next,
	}
}

func (lm *LogMiddleware) CalculateDistance(
	data types.OBUData,
) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"dist": dist,
		}).Info("Calculating distance")
	}(time.Now())
	// dist, err = lm.next.CalculateDistance(data)
	// return
	return lm.next.CalculateDistance(data)
}
