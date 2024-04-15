package main

import (
	"github.com/sirupsen/logrus"
	"github.com/swarajroy/toll_calculator/distance_calculator/consumer"
	"github.com/swarajroy/toll_calculator/distance_calculator/middleware"
	"github.com/swarajroy/toll_calculator/distance_calculator/service"
)

const (
	topic = "obu-events"
)

func main() {

	svc := service.NewCalculatoServicer()
	lm := middleware.NewLogMiddleware(svc)
	kc, err := consumer.NewDataConsumer(topic, lm)
	if err != nil {
		logrus.Fatal("kafka consumer creation errored")
	}
	kc.Start()
}
