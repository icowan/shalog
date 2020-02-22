package link

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
		"Apply": ems[:1],
	})

	r := mux.NewRouter()

	r.Handle("/link/apply", kithttp.NewServer(
		eps.ApplyEndpoint,
		decodeApplyRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodPost)

	return r
}

func decodeApplyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req applyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, ErrParams.Error())
	}

	return req, nil
}
