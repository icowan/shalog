package home

import (
	"context"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/blog/src/repository"
	"github.com/nsini/blog/src/templates"
	"net/http"
)

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	//ctx := context.Background()
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	index := kithttp.NewServer(
		makeIndexEndpoint(svc),
		decodeIndexRequest,
		encodeIndexResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/", index).Methods("GET")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		encodeError(context.Background(), repository.PostNotFound, writer)
	}).GetHandler()

	return r
}

func decodeIndexRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeIndexResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	ctx = context.WithValue(ctx, "method", "index")

	resp := response.(indexResponse)

	return templates.RenderHtml(ctx, w, resp.Data)
}

type errorer interface {
	error() error
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	switch err {
	case repository.PostNotFound:
		w.WriteHeader(http.StatusNotFound)
		ctx = context.WithValue(ctx, "method", "404")
		_ = templates.RenderHtml(ctx, w, map[string]interface{}{
			"message": err.Error(),
		})
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, _ = w.Write([]byte(err.Error()))
}
