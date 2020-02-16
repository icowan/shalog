package api

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"net/textproto"
	"strconv"
)

type PostMethod string

const (
	PostCreate     PostMethod = "metaWeblog.newPost"
	GetCategories  PostMethod = "metaWeblog.getCategories"
	NewMediaObject PostMethod = "metaWeblog.newMediaObject"
	GetPost        PostMethod = "metaWeblog.getPost"
	EditPost       PostMethod = "metaWeblog.editPost"
	GetUsersBlogs  PostMethod = "blogger.getUsersBlogs"
)

func (c PostMethod) String() string {
	return string(c)
}

var NoPermission = errors.New("not permission!")

type imageFile struct {
	Header   textproto.MIMEHeader
	Filename string
	Size     int64
	Body     []byte
}

type uploadImageRequest struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	Image    imageFile `json:"data"`
}

type imageResponse struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Filename  string `json:"filename"`
	Storename string `json:"storename"`
	Size      int64  `json:"size"`
	Path      string `json:"path"`
	Hash      string `json:"hash"`
	Timestamp int64  `json:"timestamp"`
	Url       string `json:"url"`
}

func makeUploadImageEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(uploadImageRequest)
		return s.UploadImage(ctx, req)
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		var err error
		var resp interface{}

		switch PostMethod(req.MethodName) {
		case GetUsersBlogs:
			{
				resp = &getUsersBlogsResponse{
					Params: params{
						Param: param{
							Value: value{
								Array: array{
									Data: data{
										Value: dataValue{
											Struct: valStruct{
												Member: []member{
													{Name: "isAdmin", Value: memberValue{String: "1"}},
													{Name: "url", Value: memberValue{String: "http://localhost:8080"}},
													{Name: "blogid", Value: memberValue{String: "1"}},
													{Name: "blogName", Value: memberValue{String: "nsini"}},
												},
											},
										},
									},
								},
							},
						},
					},
				}
			}
		case GetCategories:
			return s.GetCategories(ctx)
		case PostCreate:
			rs, err := s.Post(ctx, req)
			return rs, err
		case NewMediaObject:
			resp, err = s.MediaObject(ctx, req)
		case GetPost:
			{
				postId, _ := strconv.Atoi(req.Params.Param[0].Value.String)
				return s.GetPost(ctx, int64(postId))
			}
		case EditPost:
			postId, _ := strconv.ParseInt(req.Params.Param[0].Value.String, 10, 64)
			return s.EditPost(ctx, postId, req)
		}

		return resp, err
	}
}
