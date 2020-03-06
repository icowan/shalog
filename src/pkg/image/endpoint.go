package image

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/blog/src/encode"
	"github.com/pkg/errors"
	"mime/multipart"
)

type Endpoints struct {
	ListEndpoint   endpoint.Endpoint
	UploadEndpoint endpoint.Endpoint
}

type (
	listImageRequest struct {
		pageSize int
		offset   int
	}
	uploadRequest struct {
		Files []*multipart.FileHeader
	}

	imageResponse struct {
		Id        int64  `json:"id"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		Filename  string `json:"filename"`
		Storename string `json:"storename"`
		Size      string `json:"size"`
		Path      string `json:"path"`
		Hash      string `json:"hash"`
		Timestamp int64  `json:"timestamp"`
		Url       string `json:"url"`
	}
)

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		ListEndpoint: func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(listImageRequest)
			res, count, err := s.List(ctx, req.pageSize, req.offset)
			return encode.Response{
				Data: map[string]interface{}{
					"list":     res,
					"count":    count,
					"pageSize": req.pageSize,
					"offset":   req.offset,
				}, Error: err,
			}, err
		},
		UploadEndpoint: func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(uploadRequest)
			var resData []*imageResponse
			for _, f := range req.Files {
				res, err := s.UploadMedia(ctx, f)
				if err != nil {
					err = errors.Wrap(err, "s.UploadMedia")
					continue
				}
				resData = append(resData, res)
			}
			return encode.Response{
				Data:  resData,
				Error: err,
			}, err
		},
	}

	for _, m := range mdw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mdw["Upload"] {
		eps.UploadEndpoint = m(eps.UploadEndpoint)
	}

	return eps
}
