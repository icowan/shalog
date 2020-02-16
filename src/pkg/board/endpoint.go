package board

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type indexResponse struct {
	Data map[string]interface{} `json:"data,omitempty"`
	Err  error                  `json:"error,omitempty"`
}

func makeBoardEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		rs, err := s.Board(ctx)
		return indexResponse{rs, err}, err
	}
}
