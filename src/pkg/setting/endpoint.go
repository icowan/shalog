package setting

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/shalog/src/encode"
	"github.com/icowan/shalog/src/repository"
	"github.com/pkg/errors"
	"net/textproto"
)

type imageFile struct {
	Header   textproto.MIMEHeader
	Filename string
	Size     int64
	Body     []byte
}

type Endpoints struct {
	GetEndpoint         endpoint.Endpoint
	PostEndpoint        endpoint.Endpoint
	DeleteEndpoint      endpoint.Endpoint
	PutEndpoint         endpoint.Endpoint
	ListEndpoint        endpoint.Endpoint
	UploadImageEndpoint endpoint.Endpoint
	UpdateEndpoint      endpoint.Endpoint
	WechatMenuEndpoint  endpoint.Endpoint
}

type (
	settingRequest struct {
		Key         repository.SettingKey `json:"key"`
		Value       string                `json:"value"`
		Description string                `json:"description"`
	}
	imageRequest struct {
		Key       repository.SettingKey `json:"key"`
		ImageUrl  string                `json:"image_url"`
		ImageFile *imageFile            `json:"image_file"`
	}

	settingsRequest map[string]string
)

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		GetEndpoint:         makeGetEndpoint(s),
		PostEndpoint:        makePostEndpoint(s),
		DeleteEndpoint:      makeDeleteEndpoint(s),
		PutEndpoint:         makePutEndpoint(s),
		ListEndpoint:        makeListEndpoint(s),
		UploadImageEndpoint: makeUploadImageEndpoint(s),
		UpdateEndpoint:      makeUpdateEndpoint(s),
		WechatMenuEndpoint:  makeWechatMenuEndpoint(s),
	}

	for _, m := range mdw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}
	for _, m := range mdw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mdw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mdw["Put"] {
		eps.PutEndpoint = m(eps.PutEndpoint)
	}
	for _, m := range mdw["Update"] {
		eps.UpdateEndpoint = m(eps.UpdateEndpoint)
	}
	for _, m := range mdw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mdw["UploadImage"] {
		eps.UploadImageEndpoint = m(eps.UploadImageEndpoint)
	}
	for _, m := range mdw["WechatMenu"] {
		eps.WechatMenuEndpoint = m(eps.WechatMenuEndpoint)
	}

	return eps
}

func makeWechatMenuEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		err = s.WechatMenu(ctx)
		return encode.Response{
			Error: err,
		}, err
	}
}
func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(settingsRequest)
		for k, v := range req {
			if e := s.Put(ctx, repository.SettingKey(k), v, ""); e != nil {
				err = errors.Wrap(err, e.Error())
			}
		}
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(settingRequest)
		res, err := s.Get(ctx, req.Key)
		return encode.Response{Data: res, Error: err}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(settingRequest)
		err = s.Post(ctx, req.Key, req.Value, req.Description)
		return encode.Response{Error: err}, err
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(settingRequest)
		err = s.Delete(ctx, req.Key)
		return encode.Response{Error: err}, err
	}
}

func makePutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(settingRequest)
		err = s.Put(ctx, req.Key, req.Value, req.Description)
		return encode.Response{Error: err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := s.List(ctx)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeUploadImageEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(imageRequest)
		imgUrl, err := s.UploadImage(ctx, req.Key, req.ImageUrl, req.ImageFile)
		return encode.Response{Data: imgUrl, Error: err}, err
	}
}
