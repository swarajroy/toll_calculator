package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
	"github.com/swarajroy/toll_calculator/types"
)

const (
	topic = "obu-events"
)

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
}

func main() {
	log.Println("starting data receiver")
	dr := NewDataReceiver()
	http.HandleFunc("/ws", dr.handleWS)
	http.ListenAndServe(":30000", nil)
	log.Println("data receiver exited")

	time.Sleep(time.Second * 60)
}

func produceToKafka() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "topic"
	for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	log.Println("enter handleWS")
	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn
	go dr.wsReceiveLoop()
	log.Println("exit handleWS")
}

func (dr *DataReceiver) wsReceiveLoop() {
	log.Println("enter wsReceiveLoop")
	log.Println("New OBU client Connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error ", err)
			continue
		}
		log.Printf("received obu data from [%d]", data.OBUID)
		//dr.msgch <- data
	}
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.OBUData, 8),
	}
}
