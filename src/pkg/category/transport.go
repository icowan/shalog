/**
 * @Time: 2020/2/29 17:32
 * @Author: solacowa@gmail.com
 * @File: transport
 * @Software: GoLand
 */

package category

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/blog/src/encode"
	"github.com/icowan/blog/src/middleware"
	"net/http"
)

func MakeHTTPHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeJsonError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
	}

	ems := []endpoint.Middleware{
		middleware.LoginMiddleware(logger), // 0
	}

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"List": ems,
	})

	r := mux.NewRouter()

	r.Handle("/category/list", kithttp.NewServer(
		eps.ListEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}
