/**
 * @Time: 2020/2/29 17:32
 * @Author: solacowa@gmail.com
 * @File: transport
 * @Software: GoLand
 */

package category

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/blog/src/encode"
	"github.com/icowan/blog/src/middleware"
	"github.com/pkg/errors"
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
		"Post": ems,
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
	r.Handle("/category/new", kithttp.NewServer(
		eps.PostEndpoint,
		decodePostRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/category/{id:[0-9]+}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDeleteRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodDelete)
	r.Handle("/category/{id:[0-9]+}", kithttp.NewServer(
		eps.PutEndpoint,
		decodePutRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodPut)

	return r
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, ErrCategoryParams.Error())
	}

	if req.Name == "" {
		return nil, ErrCategoryParamName
	}

	return req, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req categoryRequest

	vars := mux.Vars(r)
	cid, ok := vars["id"]
	if !ok {
		return nil, ErrCategoryParams
	}
	id, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, ErrCategoryParams.Error())
	}
	req.Id = id
	return req, nil
}

func decodePutRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req categoryRequest

	vars := mux.Vars(r)
	cid, ok := vars["id"]
	if !ok {
		return nil, ErrCategoryParams
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrCategoryParams
	}

	id, err := strconv.ParseInt(cid, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, ErrCategoryParams.Error())
	}
	req.Id = id

	return req, nil
}
