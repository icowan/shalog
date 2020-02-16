package post

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/blog/src/encode"
	"github.com/icowan/blog/src/repository"
	"github.com/icowan/blog/src/repository/types"
	"github.com/icowan/blog/src/templates"
	"golang.org/x/time/rate"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var errBadRoute = errors.New("bad route")

const rateBucketNum = 6

func MakeHandler(ps Service, logger kitlog.Logger, repository repository.Repository) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeError),
	}

	ems := []endpoint.Middleware{
		UpdatePostReadNum(logger, repository),
		newTokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)),
	}

	mw := map[string][]endpoint.Middleware{
		"Get":     ems[:1],
		"Awesome": ems[1:],
		"Search":  ems[1:],
	}

	eps := NewEndpoint(ps, mw)

	r := mux.NewRouter()
	r.Handle("/post", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encodeListResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/post/popular", kithttp.NewServer(
		eps.PopularEndpoint,
		decodePopularRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/post/{id:[0-9]+}", kithttp.NewServer(
		eps.GetEndpoint,
		decodeDetailRequest,
		encodeDetailResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/post/{id:[0-9]+}", kithttp.NewServer(
		eps.AwesomeEndpoint,
		decodeAwesomeRequest,
		encode.EncodeJsonResponse,
		opts...,
	)).Methods(http.MethodPut)
	r.Handle("/search", kithttp.NewServer(
		eps.SearchEndpoint,
		decodeSearchRequest,
		encodeSearchResponse,
		opts...,
	)).Methods(http.MethodGet)
	return r
}

func decodeSearchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	keyword := r.URL.Query().Get("keyword")
	tag := r.URL.Query().Get("tag")
	category := r.URL.Query().Get("category")
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

	return searchRequest{
		Keyword:  keyword,
		Tag:      tag,
		Category: category,
		Offset:   pageOffset,
		PageSize: pageSize,
	}, nil
}

func decodeAwesomeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	return postRequest{Id: postId}, nil
}

func decodePopularRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return popularRequest{}, nil
}

func decodeListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	size := r.URL.Query().Get("pageSize")
	order := r.URL.Query().Get("order")
	by := r.URL.Query().Get("by")
	offset := r.URL.Query().Get("offset")
	action, _ := strconv.Atoi(r.URL.Query().Get("action"))

	if size == "" {
		size = "10"
	}
	if order == "" {
		order = "desc"
	}
	if by == "" {
		by = "push_time"
	}
	if offset == "" {
		offset = "0"
	}
	if action < 1 {
		action = 1
	}

	pageSize, _ := strconv.Atoi(size)
	pageOffset, _ := strconv.Atoi(offset)
	return listRequest{
		pageSize: pageSize,
		order:    order,
		by:       by,
		offset:   pageOffset,
		action:   action,
	}, nil
}

func decodeDetailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	postId, err := strconv.Atoi(id)
	if err != nil {
		return nil, errBadRoute
	}
	return postRequest{
		Id: int64(postId),
	}, nil
}

func encodeDetailResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(encode.Errorer); ok && e.Error() != nil {
		encode.EncodeError(ctx, e.Error(), w)
		return nil
	}

	ctx = context.WithValue(ctx, "method", "info")

	resp := response.(postResponse)

	return templates.RenderHtml(ctx, w, resp.Data)
}

func encodeListResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(encode.Errorer); ok && e.Error() != nil {
		encode.EncodeError(ctx, e.Error(), w)
		return nil
	}

	ctx = context.WithValue(ctx, "method", "list")

	resp := response.(listResponse)

	other := resp.Data["other"].(map[string]interface{})

	return templates.RenderHtml(ctx, w, map[string]interface{}{
		"list":      resp.Data["post"],
		"tags":      other["tags"],
		"populars":  other["populars"],
		"total":     strconv.Itoa(int(resp.Count)),
		"paginator": postPaginator("post", int(resp.Count), resp.Paginator.PageSize, resp.Paginator.Offset, ""),
	})
}

func searchReplace(data, keyword string) string {
	return strings.ReplaceAll(data, keyword, fmt.Sprintf(`<b color="red">%s<b>`, keyword))
}

func encodeSearchResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(encode.Errorer); ok && e.Error() != nil {
		encode.EncodeError(ctx, e.Error(), w)
		return nil
	}

	ctx = context.WithValue(ctx, "method", "search")

	resp := response.(listResponse)

	var rs []map[string]interface{}
	for _, val := range resp.Data["post"].([]*types.Post) {
		rs = append(rs, map[string]interface{}{
			"id":         strconv.FormatUint(uint64(val.ID), 10),
			"title":      val.Title,
			"desc":       val.Description,
			"publish_at": val.PushTime.Format("2006/01/02 15:04:05"),
			"comment":    val.Reviews,
			"author":     val.User.Username,
			"tags":       val.Tags,
		})
	}

	return templates.RenderHtml(ctx, w, map[string]interface{}{
		"list":    rs,
		"keyword": resp.Data["keyword"],
		"total":   strconv.Itoa(int(resp.Count)),
		"paginator": postPaginator("search", int(resp.Count), resp.Paginator.PageSize, resp.Paginator.Offset,
			fmt.Sprintf("&keyword=%s&tag=%s", resp.Data["keyword"], resp.Data["tag"])),
	})
}

func postPaginator(path string, count, pageSize, offset int, other string) string {
	var res []string
	var prev, next int
	prev = offset - pageSize
	next = offset + pageSize
	if offset-pageSize < 0 {
		prev = 0
	}
	if offset+pageSize > count {
		next = offset
	}
	res = append(res, fmt.Sprintf(`<a href="/%s?pageSize=10&offset=%d%s">上一页</a>&nbsp;`, path, prev, other))

	length := math.Ceil(float64(count) / float64(pageSize))
	for i := 1; i <= int(length); i++ {
		os := (i - 1) * 10
		if offset == os {
			res = append(res, fmt.Sprintf(`<b>%d</b>`, i))
			continue
		}
		res = append(res, fmt.Sprintf(`<a href="/%s?pageSize=10&offset=%d%s">%d</a>&nbsp;`, path, os, other, i))
	}
	res = append(res, fmt.Sprintf(`<a href="/%s?pageSize=10&offset=%d%s">下一页</a>&nbsp;`, path, next, other))
	return strings.Join(res, "\n")
}
