package account

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/blog/src/encode"
	"github.com/icowan/blog/src/middleware"
	"github.com/icowan/blog/src/repository"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	"net/http"
	"strings"
	"time"
)

const rateBucketNum = 2

func MakeHTTPHandler(s Service, logger kitlog.Logger, repository repository.Repository) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeJsonError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
	}

	ems := []endpoint.Middleware{
		middleware.TokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)), // 1
		middleware.LoginMiddleware(logger), // 0
	}

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Get":    ems[1:],
		"Put":    ems[1:],
		"Logout": ems[1:],
		"Login":  ems[:1],
	})

	r := mux.NewRouter()

	r.Handle("/account/info", kithttp.NewServer(
		eps.GetEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	r.Handle("/account/login", kithttp.NewServer(
		eps.LoginEndpoint,
		decodeLoginRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodPost)

	r.Handle("/account/logout", kithttp.NewServer(
		eps.LogoutEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return
		},
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	r.Handle("/account/put", kithttp.NewServer(
		eps.PutEndpoint,
		decodePutRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodPut)

	return r
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, ErrBodyParams.Error())
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Password == "" {
		return nil, ErrParamsNotNull
	}

	return req, nil
}

func decodePutRequest(_ context.Context, r *http.Request) (response interface{}, err error) {
	var req putRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, ErrBodyParams.Error())
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.Email = strings.TrimSpace(req.Email)

	// todo email校验
	// todo 用户名及密码校验

	return req, nil
}
