/**
 * @Time: 2020/5/1 21:54
 * @Author: solacowa@gmail.com
 * @File: http
 * @Software: GoLand
 */

package wechat

import (
	"context"
	"github.com/icowan/shalog/src/encode"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHTTPHandler(s Service, opts ...kithttp.ServerOption) http.Handler {

	ems := []endpoint.Middleware{}

	//s = NewLoggingServer(logger, s)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Callback": ems,
	})

	r := mux.NewRouter()

	r.Handle("/wechat/callback", kithttp.NewServer(
		eps.CallbackEndpoint,
		decodeCallbackRequest,
		encodeCallbackResponse,
		opts...,
	))

	return r
}

func decodeCallbackRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return callbackRequest{r: r}, nil
}

func encodeCallbackResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(encode.Response)
	data := resp.Data.(callbackResponse)
	data.server.ServeHTTP(w, data.r, nil)
	return nil
}
