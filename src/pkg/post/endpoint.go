package post

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/blog/src/encode"
)

type popularRequest struct {
}

type popularResponse struct {
	Data []map[string]interface{} `json:"data,omitempty"`
	Err  error                    `json:"error,omitempty"`
}

type postRequest struct {
	Id int64
}

type listRequest struct {
	order, by                string
	action, pageSize, offset int
}

type postResponse struct {
	Data map[string]interface{} `json:"data,omitempty"`
	Err  error                  `json:"error,omitempty"`
}

type paginator struct {
	By       string `json:"by,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

type listResponse struct {
	Data      map[string]interface{} `json:"data,omitempty"`
	Count     int64                  `json:"count,omitempty"`
	Paginator paginator              `json:"paginator,omitempty"`
	Err       error                  `json:"error,omitempty"`
}

func (r listResponse) error() error { return r.Err }

type searchRequest struct {
	Keyword    string
	Tag        string
	TagId      int64
	CategoryId int64
	Category   string
	Offset     int
	PageSize   int
}

type Endpoints struct {
	GetEndpoint     endpoint.Endpoint
	ListEndpoint    endpoint.Endpoint
	PopularEndpoint endpoint.Endpoint
	AwesomeEndpoint endpoint.Endpoint
	SearchEndpoint  endpoint.Endpoint
}

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		GetEndpoint:     makeGetEndpoint(s),
		ListEndpoint:    makeListEndpoint(s),
		PopularEndpoint: makePopularEndpoint(s),
		AwesomeEndpoint: makeAwesomeEndpoint(s),
		SearchEndpoint:  makeSearchEndpoint(s),
	}

	for _, m := range mdw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}
	for _, m := range mdw["Awesome"] {
		eps.AwesomeEndpoint = m(eps.AwesomeEndpoint)
	}
	for _, m := range mdw["Search"] {
		eps.SearchEndpoint = m(eps.SearchEndpoint)
	}

	return eps
}

func makeSearchEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(searchRequest)
		posts, total, err := s.Search(ctx, req.Keyword, req.Tag, req.CategoryId, req.Offset, req.PageSize)
		return listResponse{
			Data: map[string]interface{}{
				"post":    posts,
				"keyword": req.Keyword,
				"tag":     req.Tag,
			},
			Count: total,
			Paginator: paginator{
				Offset:   req.Offset,
				PageSize: req.PageSize,
			},
			Err: err,
		}, err
	}
}

func makeAwesomeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		err = s.Awesome(ctx, req.Id)
		return encode.Response{Error: err}, err
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		rs, err := s.Get(ctx, req.Id)
		return postResponse{rs, err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		rs, count, other, err := s.List(ctx, req.order, req.by, req.action, req.pageSize, req.offset)
		return listResponse{
			Data: map[string]interface{}{
				"post":  rs,
				"other": other,
			},
			Count: count,
			Paginator: paginator{
				By:       req.by,
				Offset:   req.offset,
				PageSize: req.pageSize,
			},
			Err: err,
		}, err
	}
}

func makePopularEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//req := request.(popularRequest)
		rs, err := s.Popular(ctx)
		return popularResponse{rs, err}, err
	}
}
