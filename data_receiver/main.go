package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/swarajroy/toll_calculator/types"
)

type DataReceiver struct {
	msg  chan types.OBUData
	conn *websocket.Conn
}

func main() {
	fmt.Println("Data receiver working fine")
}
