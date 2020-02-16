package api

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/blog/src/config"
	"github.com/icowan/blog/src/repository"
	"io/ioutil"
	"net/http"
)

//var errBadRoute = errors.New("bad route")
var ErrInvalidArgument = errors.New("invalid argument")

type endpoints struct {
	PostEndpoint        endpoint.Endpoint
	UploadImageEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger kitlog.Logger, repository repository.Repository, cf *config.Config) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeXmlError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
	}

	eps := endpoints{
		PostEndpoint:        makePostEndpoint(svc),
		UploadImageEndpoint: makeUploadImageEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		checkAuthMiddleware(logger, repository, cf.GetString("server", "app_key")),
	}

	ems2 := []endpoint.Middleware{
		imageCheckAuthMiddleware(logger, repository, cf.GetString("server", "app_key")),
	}

	mw := map[string][]endpoint.Middleware{
		"Post":   ems,
		"Upload": ems2,
	}

	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}

	for _, m := range mw["Upload"] {
		eps.UploadImageEndpoint = m(eps.UploadImageEndpoint)
	}

	post := kithttp.NewServer(
		eps.PostEndpoint,
		decodePostRequest,
		encodeResponse,
		opts...,
	)

	uploadImage := kithttp.NewServer(
		eps.UploadImageEndpoint,
		decodeUploadImageRequest,
		encodeJsonResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/api/post/metaweblog", post).Methods("POST")
	r.Handle("/api/upload-image", uploadImage).Methods("POST")
	return r
}

func decodeUploadImageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	username := r.Header.Get("Username")
	password := r.Header.Get("Password")

	if username == "" && password == "" {
		username = r.Form.Get("username")
		password = r.Form.Get("password")
	}

	reader, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}

	form, err := reader.ReadForm(32 << 10)
	if err != nil {
		return nil, err
	}

	if form.File == nil {
		return nil, errors.New("文件不存在")
	}

	imageFiles := form.File["markdown-image"]

	image := imageFile{}

	for _, v := range imageFiles {
		f, _ := v.Open()
		buf, err := ioutil.ReadAll(f)

		if err != nil {
			return nil, err
		}

		image.Size = v.Size
		image.Filename = v.Filename
		image.Header = v.Header
		image.Body = buf
		break
	}

	return uploadImageRequest{
		Username: username,
		Password: password,
		Image:    image,
	}, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req postRequest

	if err = xml.Unmarshal(b, &req); err != nil {
		return nil, err
	}
	switch req.MethodName {
	case NewMediaObject.String():
		{
			var req postRequest
			if err = xml.Unmarshal(b, &req); err != nil {
				return nil, err
			}
			return req, nil
		}
	}

	return req, nil
}

func encodeJsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeJsonError(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeXmlError(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	return xml.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeJsonError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrorASD:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   err.Error(),
	})
}

func encodeXmlError(ctx context.Context, err error, w http.ResponseWriter) {

	type faultStruct struct {
		Text   string    `xml:",chardata"`
		Struct valStruct `xml:"struct"`
	}

	type fault struct {
		Text  string      `xml:",chardata"`
		Value faultStruct `xml:"value"`
	}

	type errorResponse struct {
		XMLName xml.Name `xml:"methodResponse"`
		Text    string   `xml:",chardata"`
		Fault   fault    `xml:"fault"`
	}

	var faultCode string

	switch err {
	case repository.PostNotFound:
		faultCode = "404"
	case ErrInvalidArgument:
		faultCode = "401"
	case NoPermission:
		faultCode = "403"
	default:
		faultCode = "500"
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_ = xml.NewEncoder(w).Encode(errorResponse{
		Fault: fault{
			Value: faultStruct{
				Struct: valStruct{
					Member: []member{
						{Name: "faultString", Value: memberValue{String: faultCode}},
						{Name: "faultCode", Value: memberValue{String: err.Error()}},
					},
				},
			},
		},
	})
}
