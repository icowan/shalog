/**
 * @Time: 2020/3/1 16:04
 * @Author: solacowa@gmail.com
 * @File: transport
 * @Software: GoLand
 */

package tag

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/shalom/src/encode"
	"github.com/icowan/shalom/src/middleware"
	"net/http"
	"strconv"
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
		"All":  ems,
		"List": ems,
		"Post": ems,
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
	r.Handle("/tag/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/tag/new", kithttp.NewServer(
		eps.PostEndpoint,
		decodePostRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/tag/{id:[0-9]+}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDeleteRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodDelete)
	r.Handle("/tag/{id:[0-9]+}", kithttp.NewServer(
		eps.PutEndpoint,
		decodePutRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodPut)

	return r
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	tid, ok := vars["id"]
	if !ok {
		return nil, ErrTagParams
	}

	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	return tagRequest{Id: id}, err
}

func decodePutRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	tid, ok := vars["id"]
	if !ok {
		return nil, ErrTagParams
	}

	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}

	var req tagRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	req.Id = id
	return req, err
}

func decodePostRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req tagRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	keyword := r.URL.Query().Get("keyword")
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

	//pageOffset = (pageOffset - 1) * pageSize

	return tagRequest{
		offset:   pageOffset,
		pageSize: pageSize,
		tagName:  keyword,
	}, nil
}
