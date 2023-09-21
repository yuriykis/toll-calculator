package main

import (
	"github.com/sirupsen/logrus"
	"github.com/yuriykis/tolling/aggregator/client"
)

const (
	kafkaTopic         = "obu-data"
	aggregatorEndpoint = ":3001"
)

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	// httpClient := client.NewHTTPClient(aggregatorEndpoint)
	grpcClient, err := client.NewGRPCClient(aggregatorEndpoint)
	if err != nil {
		logrus.Fatal(err)
	}
	kafkaConsumer, err := NewKafkaConsumer(
		kafkaTopic,
		svc,
		grpcClient,
	)
	if err != nil {
		logrus.Fatal(err)
	}
	kafkaConsumer.Start()
}
