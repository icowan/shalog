package post

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/blog/src/encode"
	"github.com/icowan/blog/src/repository"
)

type popularRequest struct {
}

type popularResponse struct {
	Data []map[string]interface{} `json:"data,omitempty"`
	Err  error                    `json:"error,omitempty"`
}

type postResponse struct {
	Data map[string]interface{} `json:"data,omitempty"`
	Err  error                  `json:"error,omitempty"`
}

type paginator struct {
	By       string `json:"by,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
}

type listResponse struct {
	Data      map[string]interface{} `json:"data,omitempty"`
	Count     int64                  `json:"count,omitempty"`
	Paginator paginator              `json:"paginator,omitempty"`
	Err       error                  `json:"error,omitempty"`
}

func (r listResponse) error() error { return r.Err }

type Endpoints struct {
	GetEndpoint     endpoint.Endpoint
	ListEndpoint    endpoint.Endpoint
	PopularEndpoint endpoint.Endpoint
	AwesomeEndpoint endpoint.Endpoint
	SearchEndpoint  endpoint.Endpoint
	NewPostEndpoint endpoint.Endpoint
	PutPostEndpoint endpoint.Endpoint
	DeleteEndpoint  endpoint.Endpoint
	RestoreEndpoint endpoint.Endpoint
}

type (
	newPostRequest struct {
		Title       string                `json:"title"`
		Description string                `json:"description"`
		Content     string                `json:"content"`
		CategoryIds []int64               `json:"category_ids"`
		TagIds      []int64               `json:"tag_ids"`
		PostStatus  repository.PostStatus `json:"post_status"`
		Markdown    bool                  `json:"markdown"`
		ImageId     int64                 `json:"image_id"`
		Id          int64                 `json:"id"`
	}
	searchRequest struct {
		Keyword    string
		Tag        string
		TagId      int64
		CategoryId int64
		Category   string
		Offset     int
		PageSize   int
	}
	postRequest struct {
		Id int64
	}
	listRequest struct {
		order, by, category string
		pageSize, offset    int
	}
)

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		GetEndpoint:     makeGetEndpoint(s),
		ListEndpoint:    makeListEndpoint(s),
		PopularEndpoint: makePopularEndpoint(s),
		AwesomeEndpoint: makeAwesomeEndpoint(s),
		SearchEndpoint:  makeSearchEndpoint(s),
		NewPostEndpoint: makeNewPostEndpoint(s),
		PutPostEndpoint: makePutPostEndpoint(s),
		DeleteEndpoint:  makeDeleteEndpoint(s),
		RestoreEndpoint: makeRestoreEndpoint(s),
	}

	for _, m := range mdw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}
	for _, m := range mdw["Awesome"] {
		eps.AwesomeEndpoint = m(eps.AwesomeEndpoint)
	}
	for _, m := range mdw["Search"] {
		eps.SearchEndpoint = m(eps.SearchEndpoint)
	}

	// admin
	for _, m := range mdw["NewPost"] {
		eps.NewPostEndpoint = m(eps.NewPostEndpoint)
	}
	for _, m := range mdw["PutPost"] {
		eps.PutPostEndpoint = m(eps.PutPostEndpoint)
	}
	for _, m := range mdw["Delete"] {
		eps.DeleteEndpoint = m(eps.DeleteEndpoint)
	}
	for _, m := range mdw["Restore"] {
		eps.RestoreEndpoint = m(eps.RestoreEndpoint)
	}

	return eps
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		err = s.Delete(ctx, req.Id)
		return encode.Response{Error: err}, err
	}
}

func makeRestoreEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		err = s.Restore(ctx, req.Id)
		return encode.Response{Error: err}, err
	}
}

func makePutPostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(newPostRequest)
		err = s.Put(ctx, req.Id, req.Title, req.Description, req.Content, req.PostStatus,
			req.CategoryIds, req.TagIds, req.Markdown, req.ImageId)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeSearchEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(searchRequest)
		posts, total, err := s.Search(ctx, req.Keyword, req.Tag, req.CategoryId, req.Offset, req.PageSize)
		return listResponse{
			Data: map[string]interface{}{
				"post":    posts,
				"keyword": req.Keyword,
				"tag":     req.Tag,
			},
			Count: total,
			Paginator: paginator{
				Offset:   req.Offset,
				PageSize: req.PageSize,
			},
			Err: err,
		}, err
	}
}

func makeAwesomeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		err = s.Awesome(ctx, req.Id)
		return encode.Response{Error: err}, err
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		rs, err := s.Get(ctx, req.Id)
		return postResponse{rs, err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		rs, count, other, err := s.List(ctx, req.order, req.by, req.category, req.pageSize, req.offset)
		return listResponse{
			Data: map[string]interface{}{
				"post":  rs,
				"other": other,
			},
			Count: count,
			Paginator: paginator{
				By:       req.by,
				Offset:   req.offset,
				PageSize: req.pageSize,
			},
			Err: err,
		}, err
	}
}

func makePopularEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//req := request.(popularRequest)
		rs, err := s.Popular(ctx)
		return popularResponse{rs, err}, err
	}
}

func makeNewPostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(newPostRequest)
		err = s.NewPost(ctx, req.Title, req.Description, req.Content, req.PostStatus, req.CategoryIds, req.TagIds, req.Markdown, req.ImageId)
		return encode.Response{
			Error: err,
		}, err
	}
}
