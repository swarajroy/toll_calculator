package aggcleint

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/swarajroy/toll_calculator/types"
)

type HttpClient struct {
	Endpoint string
}

func NewHttpClient(endpoint string) *HttpClient {
	// var (
	// 	aggregate = endpoint + "/aggregate"
	// 	invoice   = endpoint + "/invoice"
	// )
	// logrus.Infof("aggregate endpoint = %s", aggregate)
	// logrus.Infof("invoice endpoint = %s", invoice)
	return &HttpClient{
		Endpoint: endpoint,
	}
}

func (hc *HttpClient) AggregateDistance(ctx context.Context, aggReq *types.AggregatorDistanceRequest) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", hc.Endpoint+"/aggregate", bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responsed with a status code of %d", resp.StatusCode)
	}
	return nil
}

func (hc *HttpClient) GetInvoice(ctx context.Context, invreq *types.GetInvoiceRequest) (*types.Invoice, error) {

	var (
		request = hc.Endpoint + "/invoice?obuID=" + strconv.Itoa(int(invreq.ObuID))
	)

	logrus.Infof("request -> %s", request)

	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the service responsed with a status code of %d", resp.StatusCode)
	}

	var inv types.Invoice
	if err = json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, err
	}

	return &inv, nil
}
