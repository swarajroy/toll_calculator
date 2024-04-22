package aggcleint

import (
	"context"

	"github.com/swarajroy/toll_calculator/types"
)

type Client interface {
	AggregateDistance(context.Context, *types.AggregatorDistanceRequest) error
}
