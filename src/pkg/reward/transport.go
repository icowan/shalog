/**
 * @Time: 2019-07-20 16:38
 * @Author: solacowa@gmail.com
 * @File: transport
 * @Software: GoLand
 */

package reward

import (
	"context"
	"errors"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/blog/src/encode"
	"github.com/nsini/blog/src/templates"
	"net/http"
)

var errBadRoute = errors.New("bad route")

func MakeHandler(logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
	}

	r := mux.NewRouter()
	r.Handle("/reward", kithttp.NewServer(
		func(ctx context.Context, request interface{}) (response interface{}, err error) {
			return nil, nil
		},
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
			if e, ok := response.(encode.Errorer); ok && e.Error() != nil {
				encode.EncodeError(ctx, e.Error(), w)
				return nil
			}

			ctx = context.WithValue(ctx, "method", "reward")
			return templates.RenderHtml(ctx, w, nil)
		},
		opts...,
	)).Methods("GET")

	return r
}
