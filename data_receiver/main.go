package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
	"github.com/swarajroy/toll_calculator/types"
)

const (
	OBU_EVENTS = "obu-events"
)

type DataReceiver struct {
	conn *websocket.Conn
	p    *kafka.Producer
}

func main() {
	log.Println("starting data receiver")
	dr, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", dr.handleWS)
	http.ListenAndServe(":30000", nil)
	log.Println("data receiver exited")

	time.Sleep(time.Second * 60)
}

func (dr *DataReceiver) produceToKafka(data types.OBUData) error {

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Produce messages to topic (asynchronously)
	topic := OBU_EVENTS
	return dr.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(d),
	}, nil)

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
		dr.produceToKafka(data)
	}
}

func NewDataReceiver() (*DataReceiver, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}
	//defer p.Close()

	// Delivery report handler for produced messages
	// go func() {
	// 	for e := range p.Events() {
	// 		switch ev := e.(type) {
	// 		case *kafka.Message:
	// 			if ev.TopicPartition.Error != nil {
	// 				log.Fatalf("Delivery failed: %v\n", ev.TopicPartition)
	// 			} else {
	// 				log.Printf("Delivered message to %v\n", ev.TopicPartition)
	// 			}
	// 		}
	// 	}
	// }()

	return &DataReceiver{
		p: p,
	}, nil
}
