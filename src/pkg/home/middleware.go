/**
 * @Time : 2019-09-11 13:31
 * @Author : solacowa@gmail.com
 * @File : middleware
 * @Software: GoLand
 */

package home

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/icowan/shalog/src/encode"
	"github.com/icowan/shalog/src/repository"
	"net/http"
)

func page404Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router := mux.NewRouter()
		router.NotFoundHandler = router.NewRoute().HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			encode.EncodeError(context.Background(), repository.PostNotFound, writer)
		}).GetHandler()
		next.ServeHTTP(w, r)
	})
}
