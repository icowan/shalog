package link

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/blog/src/encode"
	"github.com/icowan/blog/src/repository"
)

type (
	applyRequest struct {
		Name  string `json:"name"`
		Link  string `json:"link"`
		Icon  string `json:"icon"`
		State int    `json:"state"`
	}
	linkRequest struct {
		Id int64
	}
)

type Endpoints struct {
	ApplyEndpoint  endpoint.Endpoint
	PostEndpoint   endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
	PassEndpoint   endpoint.Endpoint
}

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ApplyEndpoint:  makeApplyEndpoint(s),
		PostEndpoint:   makePostEndpoint(s),
		DeleteEndpoint: makeDeleteEndpoint(s),
		ListEndpoint:   makeListEndpoint(s),
		PassEndpoint:   makePassEndpoint(s),
	}

	for _, m := range mdw["Apply"] {
		eps.ApplyEndpoint = m(eps.ApplyEndpoint)
	}

	for _, m := range mdw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mdw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mdw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mdw["Pass"] {
		eps.PassEndpoint = m(eps.PassEndpoint)
	}

	return eps
}

func makePassEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(linkRequest)
		err = s.Pass(ctx, req.Id)
		return encode.Response{Error: err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		links, err := s.List(ctx)
		return encode.Response{
			Data:  links,
			Error: err,
		}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(linkRequest)
		err = s.Delete(ctx, req.Id)
		return encode.Response{Error: err}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(applyRequest)
		err = s.Post(ctx, req.Name, req.Link, req.Icon, repository.LinkState(req.State))
		return encode.Response{Error: err}, err
	}
}

func makeApplyEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(applyRequest)
		err = s.ApplyLink(ctx, req.Name, req.Link, req.Icon)
		return encode.Response{Error: err}, err
	}
}
