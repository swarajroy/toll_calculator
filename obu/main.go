package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/swarajroy/toll_calculator/types"
)

const (
	duration   = 1 * time.Second
	wsEndpoint = "ws://127.0.0.1:30000/ws"
)

type OBUData struct {
	OBUID int     `json:"obuId"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func getRandomOBUIDIndex(n int) int {
	return rand.Intn(20)
}

func main() {
	n := 20
	ticker := time.NewTicker(duration)
	obuIds := generateOBUIDs(n)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {

		obuData := NewOBUData(obuIds[getRandomOBUIDIndex(n)], genCoord(), genCoord())
		fmt.Printf("%+v\n", obuData)

		if err := conn.WriteJSON(obuData); err != nil {
			log.Fatal(err)
		}
		tick := <-ticker.C
		log.Println("tick = ", tick)

	}
}

func NewOBUData(obuId int, lat, long float64) types.OBUData {
	return types.OBUData{
		OBUID: obuId,
		Lat:   lat,
		Long:  long,
	}
}

func generateOBUIDs(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
