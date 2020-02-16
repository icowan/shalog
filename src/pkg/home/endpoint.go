package home

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type indexResponse struct {
	Data map[string]interface{} `json:"data,omitempty"`
	Err  error                  `json:"error,omitempty"`
}

func makeIndexEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		rs, err := s.Index(ctx)
		return indexResponse{rs, err}, err
	}
}
