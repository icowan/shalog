package post

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/blog/src/encode"
	"github.com/icowan/blog/src/middleware"
	"github.com/icowan/blog/src/repository"
	"github.com/icowan/blog/src/repository/types"
	"github.com/icowan/blog/src/templates"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var errBadRoute = errors.New("bad route")

const rateBucketNum = 6

func MakeHandler(ps Service, logger kitlog.Logger, repository repository.Repository, settings map[string]string) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerBefore(middleware.SettingsRequest(settings)),
		kithttp.ServerErrorEncoder(encode.EncodeError),
	}

	adminOpts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encode.EncodeJsonError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
	}

	ems := []endpoint.Middleware{
		UpdatePostReadNum(logger, repository),
		middleware.TokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)), // 限流
		middleware.LoginMiddleware(logger),                                                        // 是否登录
	}

	eps := NewEndpoint(ps, map[string][]endpoint.Middleware{
		"Get":     ems[:1],
		"Awesome": ems[2:2],
		"Search":  ems[2:2],
		"NewPost": ems[2:],
		"PutPost": ems[2:],
		"Delete":  ems[2:],
		"Restore": ems[2:],
		"Detail":  ems[2:],
		"Star":    ems[2:],
	})

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
	r.Handle("/post/{id:[0-9]+}/awesome", kithttp.NewServer(
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

	// 后端API
	r.Handle("/post/list", kithttp.NewServer(
		eps.AdminListEndpoint,
		decodeListRequest,
		encode.EncodeJsonResponse,
		adminOpts...,
	)).Methods(http.MethodGet)
	r.Handle("/post/new", kithttp.NewServer(
		eps.NewPostEndpoint,
		decodeNewPostRequest,
		encode.EncodeJsonResponse,
		adminOpts...,
	)).Methods(http.MethodPost)
	r.Handle("/post/{id:[0-9]+}", kithttp.NewServer(
		eps.PutPostEndpoint,
		decodePutPostRequest,
		encode.EncodeJsonResponse,
		adminOpts...,
	)).Methods(http.MethodPut)
	r.Handle("/post/{id:[0-9]+}/detail", kithttp.NewServer(
		eps.DetailEndpoint,
		decodeDetailRequest,
		encode.EncodeJsonResponse,
		adminOpts...,
	)).Methods(http.MethodGet)
	r.Handle("/post/{id:[0-9]+}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeDetailRequest,
		encode.EncodeJsonResponse,
		adminOpts...,
	)).Methods(http.MethodDelete)
	r.Handle("/post/{id:[0-9]+}/restore", kithttp.NewServer(
		eps.RestoreEndpoint,
		decodeDetailRequest,
		encode.EncodeJsonResponse,
		adminOpts...,
	)).Methods(http.MethodPut)
	r.Handle("/post/{id:[0-9]+}/star", kithttp.NewServer(
		eps.StarEndpoint,
		decodeDetailRequest,
		encode.EncodeJsonResponse,
		adminOpts...,
	)).Methods(http.MethodPut)

	return r
}

func decodePutPostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	var req newPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, ErrPostParams.Error())
	}

	req.Id = postId
	return req, nil
}

func trimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

func decodeNewPostRequest(_ context.Context, r *http.Request) (response interface{}, err error) {
	req := newPostRequest{
		PostStatus: repository.PostStatusDraft,
		Markdown:   false,
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	req.Title = trimHtml(req.Title)

	if req.Title == "" {
		return nil, ErrPostParamTitle
	}

	if len(req.Categories) < 1 {
		return nil, ErrPostParamCategories
	}

	// todo: 一堆验证......

	return req, nil
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
	category := strings.TrimSpace(r.URL.Query().Get("category"))
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))

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

	pageSize, _ := strconv.Atoi(size)
	pageOffset, _ := strconv.Atoi(offset)
	return listRequest{
		pageSize: pageSize,
		order:    order,
		by:       by,
		offset:   pageOffset,
		category: category,
		keyword:  keyword,
	}, nil
}

func decodeDetailRequest(ctx context.Context, r *http.Request) (interface{}, error) {
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
		"list":       resp.Data["post"],
		"tags":       other["tags"],
		"populars":   other["populars"],
		"categories": other["categories"],
		"total":      strconv.Itoa(int(resp.Count)),
		"paginator": postPaginator("post", int(resp.Count), resp.Paginator.PageSize, resp.Paginator.Offset,
			"&category="+other["category"].(string)),
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
