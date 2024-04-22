package main

import (
	"github.com/sirupsen/logrus"
	"github.com/swarajroy/toll_calculator/aggregator/aggcleint"
	"github.com/swarajroy/toll_calculator/distance_calculator/consumer"
	"github.com/swarajroy/toll_calculator/distance_calculator/middleware"
	"github.com/swarajroy/toll_calculator/distance_calculator/service"
)

const (
	TOPIC    = "obu-events"
	ENDPOINT = "http://localhost:3000/aggregate"
)

func main() {
	client := aggcleint.NewHttpClient(ENDPOINT)
	svc := service.NewCalculatoServicer()
	lm := middleware.NewLogMiddleware(svc)
	kc, err := consumer.NewDataConsumer(TOPIC, lm, client)
	if err != nil {
		logrus.Fatal("kafka consumer creation errored")
	}
	kc.Start()
}
