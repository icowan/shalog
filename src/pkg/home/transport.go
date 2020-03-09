package home

import (
	"context"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/shalom/src/encode"
	"github.com/icowan/shalom/src/middleware"
	"github.com/icowan/shalom/src/repository"
	"github.com/icowan/shalom/src/templates"
	"net/http"
)

func MakeHandler(svc Service, logger kitlog.Logger, settings map[string]string) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerBefore(middleware.SettingsRequest(settings)),
		kithttp.ServerErrorEncoder(encode.EncodeError),
	}

	r := mux.NewRouter()
	r.Handle("/", kithttp.NewServer(
		makeIndexEndpoint(svc),
		decodeIndexRequest,
		encodeIndexResponse,
		opts...,
	)).Methods(http.MethodGet)

	r.NotFoundHandler = r.NewRoute().HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "settings", settings)
		encode.EncodeError(ctx, repository.PostNotFound, writer)
	}).GetHandler()

	return r
}

func decodeIndexRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeIndexResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(encode.Errorer); ok && e.Error() != nil {
		encode.EncodeError(ctx, e.Error(), w)
		return nil
	}

	ctx = context.WithValue(ctx, "method", "index")

	resp := response.(indexResponse)

	return templates.RenderHtml(ctx, w, resp.Data)
}
