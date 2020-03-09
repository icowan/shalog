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
	"github.com/icowan/shalom/src/encode"
)

type Endpoints struct {
	AllEndpoint    endpoint.Endpoint
	PostEndpoint   endpoint.Endpoint
	PutEndpoint    endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
	GetEndpoint    endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
}

type tagRequest struct {
	offset, pageSize int
	tagName          string
	Name             string `json:"name"`
	Description      string `json:"description"`
	Id               int64  `json:"id"`
}

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		AllEndpoint:    makeAllEndpoint(s),
		ListEndpoint:   makeListEndpoint(s),
		PostEndpoint:   makePostEndpoint(s),
		PutEndpoint:    makePutEndpoint(s),
		DeleteEndpoint: makeDeleteEndpoint(s),
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
	for _, m := range mdw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}

	return eps
}

func makePutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(tagRequest)
		err = s.Put(ctx, req.Id, req.Name)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(tagRequest)
		err = s.Delete(ctx, req.Id)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(tagRequest)
		err = s.Post(ctx, req.Name)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := s.All(ctx)
		return encode.Response{
			Data: res,
		}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(tagRequest)
		res, count, err := s.List(ctx, req.tagName, req.pageSize, req.offset)
		return encode.Response{
			Data: map[string]interface{}{
				"list":     res,
				"count":    count,
				"pageSize": req.pageSize,
				"offset":   req.offset,
			},
			Error: err,
		}, err
	}
}
