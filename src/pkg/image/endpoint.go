package image

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/blog/src/encode"
)

type Endpoints struct {
	ListEndpoint endpoint.Endpoint
}

type listImageRequest struct {
	pageSize int
	offset   int
}

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListEndpoint: func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(listImageRequest)
			res, count, err := s.List(ctx, req.pageSize, req.offset)
			return encode.Response{
				Data: map[string]interface{}{
					"list":     res,
					"count":    count,
					"pageSize": req.pageSize,
					"offset":   req.offset,
				}, Error: err,
			}, err
		}}

	for _, m := range mdw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}

	return eps
}
