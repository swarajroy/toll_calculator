package aggcleint

import (
	"github.com/swarajroy/toll_calculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Endpoint string
	types.DistanceAggregatorClient
}

func NewGrpcClient(endpoint string) (*GrpcClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewDistanceAggregatorClient(conn)
	return &GrpcClient{
		Endpoint:                 endpoint,
		DistanceAggregatorClient: c,
	}, nil
}
