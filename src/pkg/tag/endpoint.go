/**
 * @Time: 2020/3/1 15:56
 * @Author: solacowa@gmail.com
 * @File: endpoint
 * @Software: GoLand
 */

package tag

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/blog/src/encode"
)

type Endpoints struct {
	AllEndpoint    endpoint.Endpoint
	PostEndpoint   endpoint.Endpoint
	PutEndpoint    endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
	GetEndpoint    endpoint.Endpoint
}

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		AllEndpoint:    makeAllEndpoint(s),
		PostEndpoint:   nil,
		PutEndpoint:    nil,
		DeleteEndpoint: nil,
		GetEndpoint:    nil,
	}

	for _, m := range mdw["All"] {
		eps.AllEndpoint = m(eps.AllEndpoint)
	}
	for _, m := range mdw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mdw["Put"] {
		eps.PutEndpoint = m(eps.PutEndpoint)
	}
	for _, m := range mdw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mdw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}

	return eps
}

func makeAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := s.All(ctx)
		return encode.Response{
			Data: res,
		}, err
	}
}
