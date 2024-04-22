package aggcleint

import (
	"context"

	"github.com/swarajroy/toll_calculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Endpoint string
	client   types.DistanceAggregatorClient
}

func NewGrpcClient(endpoint string) (*GrpcClient, error) {
	conn, err := grpc.Dial(":3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewDistanceAggregatorClient(conn)
	return &GrpcClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (c *GrpcClient) AggregateDistance(ctx context.Context, req *types.AggregatorDistanceRequest) error {
	_, err := c.client.AggregateDistance(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
