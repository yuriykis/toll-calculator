package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yuriykis/tolling/types"
)

type LogModdleware struct {
	next CalculatorServicer
}

func NewLogModdleware(next CalculatorServicer) CalculatorServicer {
	return &LogModdleware{
		next: next,
	}
}

func (lm *LogModdleware) CalculateDistance(data types.OBUData) (dist float64, err error) {
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
