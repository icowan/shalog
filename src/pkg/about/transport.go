package about

import (
	"context"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/blog/src/repository"
	"github.com/icowan/blog/src/templates"
	"net/http"
)

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	//ctx := context.Background()
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	about := kithttp.NewServer(
		makeAboutEndpoint(svc),
		decodeAboutRequest,
		encodeAboutResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/about", about).Methods("GET")
	return r
}

func decodeAboutRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeAboutResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	ctx = context.WithValue(ctx, "method", "about")

	return templates.RenderHtml(ctx, w, map[string]interface{}{
		"title": "about",
	})
}

type errorer interface {
	error() error
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	switch err {
	case repository.PostNotFound:
		w.WriteHeader(http.StatusNotFound)
		ctx = context.WithValue(ctx, "method", "404")
		_ = templates.RenderHtml(ctx, w, map[string]interface{}{})
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, _ = w.Write([]byte(err.Error()))
}
