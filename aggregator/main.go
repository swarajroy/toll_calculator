package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/swarajroy/toll_calculator/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpListenAddr", ":3000", "the listening address of the HTTP Server")
	grpcListenAddr := flag.String("grpcListenAddr", ":3001", "the listening address of the GRPC Server")
	flag.Parse()

	var (
		store = NewInMemoryStore()
		svc   = NewAggregator(store)
	)
	svc = NewLogMiddleware(svc)
	go makeGRPCTransport(*grpcListenAddr, svc)
	//time.Sleep(time.Second * 2)
	// c, err := aggcleint.NewGrpcClient(*grpcListenAddr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// c.AggregateDistance(context.Background(), &types.AggregatorDistanceRequest{
	// 	ObuID: 1,s
	// 	Value: 58.55,
	// 	Unix:  time.Now().UnixNano(),
	// })
	makeHTTPTransport(*httpListenAddr, svc)
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	// Make a TCP Listener
	fmt.Println("GRPC Transport running on ", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	// Make a native GRPC Server
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// Register our DistanceAggServer
	types.RegisterDistanceAggregatorServer(server, NewGRPCDistanceAggregatorServer(svc))
	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP Transport running on ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuID, err := strconv.Atoi(r.URL.Query().Get("obuID"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu id"})
			return
		}
		invoice, err := svc.Invoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "calculating the invoice for the obuID errored"})
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
