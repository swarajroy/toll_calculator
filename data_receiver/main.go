package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/swarajroy/toll_calculator/types"
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
