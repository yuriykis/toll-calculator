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
	svc = NewLogMiddleware(svc)

	kafkaConsumer, err := NewKafkaConsumer(
		kafkaTopic,
		svc,
		client.NewHTTPClient(aggregatorEndpoint),
	)
	if err != nil {
		logrus.Fatal(err)
	}
	kafkaConsumer.Start()
}
