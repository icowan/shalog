package link

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/blog/src/encode"
)

type (
	applyRequest struct {
		Name string `json:"name"`
		Link string `json:"link"`
		Icon string `json:"icon"`
	}
)

type Endpoints struct {
	ApplyEndpoint  endpoint.Endpoint
	PostEndpoint   endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
}

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ApplyEndpoint:  makeApplyEndpoint(s),
		PostEndpoint:   nil,
		DeleteEndpoint: nil,
		ListEndpoint:   nil,
	}

	for _, m := range mdw["Apply"] {
		eps.ApplyEndpoint = m(eps.ApplyEndpoint)
	}

	return eps
}

func makeApplyEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(applyRequest)
		err = s.ApplyLink(ctx, req.Name, req.Link, req.Icon)
		return encode.Response{Error: err}, err
	}
}
