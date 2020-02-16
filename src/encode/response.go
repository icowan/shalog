/**
 * @Time: 2020/2/14 13:57
 * @Author: solacowa@gmail.com
 * @File: encode
 * @Software: GoLand
 */

package encode

import (
	"context"
	"encoding/json"
	"github.com/icowan/blog/src/repository"
	"github.com/icowan/blog/src/templates"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

type Failure interface {
	Failed() error
}

type Errorer interface {
	Error() error
}

func err2code(err error) int {
	return http.StatusOK
}

type errorWrapper struct {
	Error string `json:"error"`
}

func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
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

func EncodeJsonError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusOK)
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func EncodeJsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(Failure); ok && f.Failed() != nil {
		EncodeJsonError(ctx, f.Failed(), w)
		return nil
	}
	resp := response.(Response)
	if resp.Error == nil {
		resp.Success = true
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(resp)
}
