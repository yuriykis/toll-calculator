package main

import (
	"github.com/sirupsen/logrus"
	"github.com/yuriykis/tolling/aggregator/client"
)

const (
	kafkaTopic         = "obu-data"
	aggregatorEndpoint = "http://localhost:3000"
)

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogModdleware(svc)

	kafkaConsumer, err := NewKafkaConsumer(
		kafkaTopic,
		svc,
		client.NewClient(aggregatorEndpoint),
	)
	if err != nil {
		logrus.Fatal(err)
	}
	kafkaConsumer.Start()
}
