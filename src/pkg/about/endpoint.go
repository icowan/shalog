package about

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type aboutResponse struct {
	Data map[string]interface{} `json:"data,omitempty"`
	Err  error                  `json:"error,omitempty"`
}

func makeAboutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		rs, err := s.About(ctx)
		return aboutResponse{rs, err}, err
	}
}
