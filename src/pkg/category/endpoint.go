/**
 * @Time: 2020/2/29 17:20
 * @Author: solacowa@gmail.com
 * @File: endpoint
 * @Software: GoLand
 */

package category

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/shalog/src/encode"
)

type Endpoints struct {
	ListEndpoint   endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
	PutEndpoint    endpoint.Endpoint
	PostEndpoint   endpoint.Endpoint
}

type (
	categoryRequest struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		ParentId    int64  `json:"parent_id"`
	}
)

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListEndpoint:   makeListEndpoint(s),
		DeleteEndpoint: makeDeleteEndpoint(s),
		PutEndpoint:    makePutEndpoint(s),
		PostEndpoint:   makePostEndpoint(s),
	}

	for _, m := range mdw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
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

	return eps
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		categories, err := s.List(ctx)
		return encode.Response{
			Data:  categories,
			Error: err,
		}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(categoryRequest)
		err = s.Post(ctx, req.Name, req.Description, req.ParentId)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makePutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(categoryRequest)
		err = s.Put(ctx, req.Id, req.Name, req.Description, req.ParentId)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(categoryRequest)
		err = s.Delete(ctx, req.Id)
		return encode.Response{
			Error: err,
		}, err
	}
}
