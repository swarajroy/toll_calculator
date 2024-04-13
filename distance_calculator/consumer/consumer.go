package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"github.com/swarajroy/toll_calculator/types"
)

type DataConsumer interface {
	Consume(data *types.OBUData)
	Start()
}

type KafkaDataConsumer struct {
	c         *kafka.Consumer
	isRunning bool
}

func NewDataConsumer(topic string) (DataConsumer, error) {
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
		c: c,
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
		}
		kc.Consume(data)
	}
}

func (kc *KafkaDataConsumer) Consume(data *types.OBUData) {
	logrus.WithFields(logrus.Fields{
		"obuID": data.OBUID,
	}).Info("consuming from kafka")
}

func (kc *KafkaDataConsumer) Start() {
	kc.isRunning = true
	kc.readMessageLoop()
}
