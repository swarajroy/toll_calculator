package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"github.com/swarajroy/toll_calculator/aggregator/aggcleint"
	"github.com/swarajroy/toll_calculator/distance_calculator/service"
	"github.com/swarajroy/toll_calculator/types"
)

type DataConsumer interface {
	Start()
}

type KafkaDataConsumer struct {
	c         *kafka.Consumer
	isRunning bool
	svc       service.CalculatorServicer
	client    *aggcleint.Client
}

func NewDataConsumer(topic string, svc service.CalculatorServicer, client *aggcleint.Client) (DataConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaDataConsumer{
		c:      c,
		svc:    svc,
		client: client,
	}, nil
}

func (kc *KafkaDataConsumer) readMessageLoop() {
	for kc.isRunning {
		msg, err := kc.c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		}

		// UnMarshal the []byte of msg into data of types.OBUData
		var data *types.OBUData
		if err = json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("error occured whilst unmarshalling")
			continue
		}

		distance, err := kc.svc.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("error occurred")
			continue
		}
		err = kc.client.AggregateDistance(types.NewDistance(distance, data.OBUID, time.Now().UnixNano()))
		if err != nil {
			logrus.Error("aggregate error", err.Error())
			continue
		}
	}
	fmt.Println("exit readMessageLoop")
	kc.c.Close()
}

func (kc *KafkaDataConsumer) Start() {
	fmt.Println("kafka transport for consumption started")
	kc.isRunning = true
	kc.readMessageLoop()
}
