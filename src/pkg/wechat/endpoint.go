/**
 * @Time: 2020/5/1 21:59
 * @Author: solacowa@gmail.com
 * @File: endpoint
 * @Software: GoLand
 */

package wechat

import (
	"context"
	"github.com/chanxuehong/wechat/mp/core"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/icowan/shalog/src/encode"
)

type (
	callbackRequest struct {
		r *http.Request
	}

	callbackResponse struct {
		r      *http.Request
		server *core.Server
	}
)

type Endpoints struct {
	CallbackEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CallbackEndpoint: makeCallbackEndpoint(s),
	}

	for _, m := range mdw["Callback"] {
		eps.CallbackEndpoint = m(eps.CallbackEndpoint)
	}

	return eps
}

func makeCallbackEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(callbackRequest)
		server := s.Callback(ctx)

		return encode.Response{
			Data: callbackResponse{
				r:      req.r,
				server: server,
			},
			Error: err,
		}, err
	}
}
