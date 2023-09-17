package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yuriykis/tolling/types"
)

type LogModdleware struct {
	next DataProducer
}

func NewLogModdleware(next DataProducer) *LogModdleware {
	return &LogModdleware{
		next: next,
	}
}

func (lm *LogModdleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUID,
			"lat":   data.Lat,
			"long":  data.Long,
			"took":  time.Since(start),
		}).Info("Producing data")
	}(time.Now())
	return lm.next.ProduceData(data)
}
