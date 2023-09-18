package main

import "github.com/sirupsen/logrus"

const kafkaTopic = "obu-data"

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogModdleware(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		logrus.Fatal(err)
	}
	kafkaConsumer.Start()
}
