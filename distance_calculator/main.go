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
	ENDPOINT = "http://127.0.0.1:3000/aggregate"
)

func main() {
	//httpclient := aggcleint.NewHttpClient(ENDPOINT)
	grpcClient, err := aggcleint.NewGrpcClient(ENDPOINT)

	if err != nil {
		logrus.Fatal("grpc client construction errored")
	}

	svc := service.NewCalculatoServicer()
	lm := middleware.NewLogMiddleware(svc)
	kc, err := consumer.NewDataConsumer(TOPIC, lm, grpcClient)
	if err != nil {
		logrus.Fatal("kafka consumer creation errored")
	}
	kc.Start()
}
