package producer

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/swarajroy/toll_calculator/types"
)

const (
	OBU_EVENTS = "obu-events"
)

type DataProducer interface {
	Publish(data types.OBUData) error
}

type KafkaProducer struct {
	p *kafka.Producer
}

func NewKafkaProducer() (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}
	defer p.Close()

	//Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Fatalf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &KafkaProducer{
		p: p,
	}, nil
}

func (kp *KafkaProducer) Publish(data types.OBUData) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Produce messages to topic (asynchronously)
	topic := OBU_EVENTS
	return kp.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(d),
	}, nil)
}
