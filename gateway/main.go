package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/swarajroy/toll_calculator/aggregator/aggcleint"
	"github.com/swarajroy/toll_calculator/types"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type GatewayApieHandler struct {
	client aggcleint.Client
}

func NewGatewayApiHandler(client aggcleint.Client) *GatewayApieHandler {
	return &GatewayApieHandler{
		client: client,
	}
}

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "listen address of the http server")
	flag.Parse()

	var (
		client            = aggcleint.NewHttpClient("http://127.0.0.1:3000")
		gatewayApiHandler = NewGatewayApiHandler(client)
	)

	http.HandleFunc("/invoice", makeAPIFunc(gatewayApiHandler.handleGetInvoice))
	logrus.Infof("gateway http server running on port %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func (gh *GatewayApieHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	obuID, err := strconv.Atoi((r.URL.Query().Get("obuID")))
	if err != nil {
		return err
	}
	inv, err := gh.client.GetInvoice(context.Background(), &types.GetInvoiceRequest{
		ObuID: int64(obuID),
	})
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri":  r.RequestURI,
			}).Info("REQ::")
		}(time.Now())
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
