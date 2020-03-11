package link

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/shalog/src/encode"
	"github.com/icowan/shalog/src/middleware"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const rateBucketNum = 2

var errBadRoute = errors.New("bad route")

func MakeHTTPHandler(s Service, logger kitlog.Logger) http.Handler {
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
		"Apply":  ems[:1],
		"List":   ems[:1],
		"Post":   ems[1:],
		"Pass":   ems[1:],
		"Delete": ems[1:],
		"All":    ems[1:],
	})

	r := mux.NewRouter()

	r.Handle("/link/apply", kithttp.NewServer(
		eps.ApplyEndpoint,
		decodeApplyRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/link/list", kithttp.NewServer(
		eps.ListEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	r.Handle("/link/all", kithttp.NewServer(
		eps.AllEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	r.Handle("/link/new", kithttp.NewServer(
		eps.PostEndpoint,
		decodeApplyRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodPost)

	r.Handle("/link/{id:[0-9]+}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDeleteRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodDelete)

	r.Handle("/link/{id:[0-9]+}/pass", kithttp.NewServer(
		eps.PassEndpoint,
		decodeDeleteRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodPut)

	return r
}

func decodeApplyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	b, _ := ioutil.ReadAll(r.Body)
	var req applyRequest
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	return nil, errors.Wrap(err, ErrParams.Error())
	//}
	if err := json.Unmarshal(b, &req); err != nil {
		return nil, errors.Wrap(err, ErrParams.Error())
	}

	if len(req.Link) > 1000 {
		return nil, ErrParamLinkLen
	}

	if len(req.Link) > 64 {
		return nil, ErrParamNameLen
	}

	// todo 解析link还有问题

	if _, err := url.ParseRequestURI(req.Link); err != nil {
		return nil, errors.Wrap(err, ErrParams.Error())
	}
	if len(strings.TrimSpace(req.Icon)) > 3 {
		if _, err := url.ParseRequestURI(req.Icon); err != nil {
			return nil, errors.Wrap(err, ErrParams.Error())
		}
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Link = url.QueryEscape(strings.TrimSpace(req.Link))
	req.Icon = url.QueryEscape(strings.TrimSpace(req.Icon))

	if req.Name == "" || req.Link == "" {
		return nil, ErrParams
	}

	return req, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req linkRequest

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	linkId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, errBadRoute.Error())
	}
	req.Id = linkId
	return req, nil
}
