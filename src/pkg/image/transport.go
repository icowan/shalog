package image

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/shalom/src/encode"
	"github.com/icowan/shalom/src/middleware"
	"github.com/icowan/shalom/src/repository"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
)

func MakeHTTPHandler(s Service, logger kitlog.Logger, settings map[string]string) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeJsonError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(middleware.SettingsRequest(settings)),
	}

	getImgOpts := []kithttp.ServerOption{
		kithttp.ServerBefore(middleware.SettingsRequest(settings)),
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
	}

	ems := []endpoint.Middleware{
		middleware.LoginMiddleware(logger), // 0
	}

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"List": ems,
	})

	r := mux.NewRouter()

	r.Handle("/image/list", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodGet)

	r.Handle("/image/upload", kithttp.NewServer(
		eps.UploadEndpoint,
		decodeUploadRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods("POST")

	imageDomain := settings[repository.SettingGlobalDomainImage.String()]
	u, _ := url.Parse(imageDomain)

	r.PathPrefix(u.Path + "/").Handler(kithttp.NewServer(
		eps.GetEndpoint,
		decodeGetRequest,
		encodeImageResponse,
		getImgOpts...,
	)).Methods(http.MethodGet)
	return r
}

func decodeGetRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return imageRequest{
		Url: r.URL,
	}, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
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

	return listImageRequest{
		offset:   pageOffset,
		pageSize: pageSize,
	}, nil
}

func decodeUploadRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	reader, err := r.MultipartReader()
	if err != nil {
		return nil, errors.Wrap(err, "r.MultipartReader")
	}

	form, err := reader.ReadForm(32 << 10)
	if err != nil {
		return nil, errors.Wrap(err, "reader.ReadForm")
	}

	if form.File == nil {
		return nil, errors.New("文件不存在")
	}

	return uploadRequest{Files: form.File["file"]}, nil
}

func encodeImageResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if e, ok := response.(encode.Errorer); ok && e.Error() != nil {
		encode.EncodeError(ctx, e.Error(), w)
		return err
	}

	res := response.(encode.Response)
	b := res.Data.([]byte)

	_, err = w.Write(b)
	return
}
