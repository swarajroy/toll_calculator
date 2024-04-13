package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/swarajroy/toll_calculator/distance_calculator/consumer"
)

const (
	topic = "obu-events"
)

func main() {
	fmt.Println("This is working just fine")
	kc, err := consumer.NewDataConsumer(topic)
	if err != nil {
		logrus.Fatal("kafka consumer creation errored")
	}
	kc.Start()
}
