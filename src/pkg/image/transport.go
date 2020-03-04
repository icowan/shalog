package image

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/blog/src/encode"
	"github.com/icowan/blog/src/middleware"
	"net/http"
	"strconv"
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

	r.Handle("/image/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	return r
}

func decodeListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	offset := r.URL.Query().Get("offset")
	size := r.URL.Query().Get("pageSize")

	if size == "" {
		size = "10"
	}
	if offset == "" {
		offset = "0"
	}
	pageSize, _ := strconv.Atoi(size)
	pageOffset, _ := strconv.Atoi(offset)

	return listImageRequest{
		offset:   pageOffset,
		pageSize: pageSize,
	}, nil
}
