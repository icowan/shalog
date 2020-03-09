package post

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/shalom/src/encode"
	"github.com/icowan/shalom/src/repository"
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
	GetEndpoint       endpoint.Endpoint
	ListEndpoint      endpoint.Endpoint
	PopularEndpoint   endpoint.Endpoint
	AwesomeEndpoint   endpoint.Endpoint
	SearchEndpoint    endpoint.Endpoint
	NewPostEndpoint   endpoint.Endpoint
	PutPostEndpoint   endpoint.Endpoint
	DeleteEndpoint    endpoint.Endpoint
	RestoreEndpoint   endpoint.Endpoint
	AdminListEndpoint endpoint.Endpoint
	DetailEndpoint    endpoint.Endpoint
	StarEndpoint      endpoint.Endpoint
}

type (
	newPostRequest struct {
		Title       string                `json:"title"`
		Description string                `json:"description"`
		Content     string                `json:"content"`
		CategoryIds []int64               `json:"category_ids"`
		TagIds      []int64               `json:"tag_ids"`
		Tags        []string              `json:"tags"`
		Categories  []string              `json:"categories"`
		PostStatus  repository.PostStatus `json:"post_status"`
		Markdown    bool                  `json:"is_markdown"`
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
		order, by, category, tag, keyword string
		pageSize, offset                  int
	}
)

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		GetEndpoint:       makeGetEndpoint(s),
		ListEndpoint:      makeListEndpoint(s),
		PopularEndpoint:   makePopularEndpoint(s),
		AwesomeEndpoint:   makeAwesomeEndpoint(s),
		SearchEndpoint:    makeSearchEndpoint(s),
		NewPostEndpoint:   makeNewPostEndpoint(s),
		PutPostEndpoint:   makePutPostEndpoint(s),
		DeleteEndpoint:    makeDeleteEndpoint(s),
		RestoreEndpoint:   makeRestoreEndpoint(s),
		AdminListEndpoint: makeAdminListEndpoint(s),
		DetailEndpoint:    makeDetailEndpoint(s),
		StarEndpoint:      makeStarEndpoint(s),
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
	for _, m := range mdw["Detail"] {
		eps.DetailEndpoint = m(eps.DetailEndpoint)
	}

	return eps
}

func makeStarEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		err = s.Star(ctx, req.Id)
		return encode.Response{Error: err}, err
	}
}

func makeDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		rs, err := s.Detail(ctx, req.Id)
		return encode.Response{
			Data:  rs,
			Error: err,
		}, err
	}
}

func makeAdminListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		posts, count, err := s.AdminList(ctx, req.order, "created_at", req.category, req.tag, req.pageSize, req.offset, req.keyword)
		return encode.Response{
			Error: err,
			Data: map[string]interface{}{
				"list":     posts,
				"count":    count,
				"offset":   req.offset,
				"pageSize": req.pageSize,
			}}, err
	}
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
			req.Categories, req.Tags, req.Markdown, req.ImageId)
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

// TODO: 继续下楼取快递

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		rs, count, other, err := s.List(ctx, req.order, req.by, req.category, req.pageSize, req.offset)
		populars, _ := s.Popular(ctx)
		other["populars"] = populars
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
		var res interface{}
		if req.Id > 0 {
			err = s.Put(ctx, req.Id, req.Title, req.Description, req.Content, req.PostStatus, req.Categories, req.Tags, req.Markdown, req.ImageId)
		} else {
			res, err = s.NewPost(ctx, req.Title, req.Description, req.Content, req.PostStatus, req.Categories, req.Tags, req.Markdown, req.ImageId)
		}
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}
