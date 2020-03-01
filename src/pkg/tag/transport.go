/**
 * @Time: 2020/3/1 16:04
 * @Author: solacowa@gmail.com
 * @File: transport
 * @Software: GoLand
 */

package tag

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/blog/src/encode"
	"github.com/icowan/blog/src/middleware"
	"net/http"
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeJsonError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
	}

	ems := []endpoint.Middleware{
		middleware.LoginMiddleware(logger), // 0
	}

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"All": ems,
	})

	r := mux.NewRouter()

	r.Handle("/tag/all", kithttp.NewServer(
		eps.AllEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}
