package home

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/blog/src/config"
	"github.com/icowan/blog/src/pkg/post"
	"github.com/icowan/blog/src/repository"
	"strconv"
)

type Service interface {
	Index(ctx context.Context) (rs map[string]interface{}, err error)
	//ApplyLink(ctx context.Context, name, link, icon string) (err error)

}

type service struct {
	logger     log.Logger
	config     *config.Config
	repository repository.Repository
}

func NewService(logger log.Logger, config *config.Config, repository repository.Repository) Service {
	return &service{
		logger:     logger,
		config:     config,
		repository: repository,
	}
}

func (c *service) Index(ctx context.Context) (rs map[string]interface{}, err error) {
	stars, err := c.repository.Post().Stars()
	if err != nil {
		_ = level.Error(c.logger).Log("Post", "Stars", "err", err.Error())
		return
	}

	var starsData []map[string]interface{}

	for _, v := range stars {
		var imgUrl string
		if len(v.Images) > 0 {
			imgUrl = c.config.GetString("server", "image_domain") + "/" + v.Images[0].ImagePath
		}
		starsData = append(starsData, map[string]interface{}{
			"content":    v.Content,
			"title":      v.Title,
			"publish_at": v.PushTime.Format("2006/01/02 15:04:05"),
			"updated_at": v.UpdatedAt,
			"author":     v.User.Username,
			"comment":    v.Reviews,
			"image_url":  imgUrl,
			"desc":       v.Description,
			"id":         strconv.Itoa(int(v.ID)),
		})
	}

	list, err := c.repository.Post().Index()
	if err != nil {
		_ = level.Error(c.logger).Log("Post", "Index", "err", err.Error())
		return
	}

	var posts []map[string]interface{}
	for _, v := range list {
		var imgUrl string
		if len(v.Images) > 0 {
			imgUrl = c.config.GetString("server", "image_domain") + "/" + v.Images[0].ImagePath
		}
		posts = append(posts, map[string]interface{}{
			"content":    v.Content,
			"title":      v.Title,
			"publish_at": v.PushTime.Format("2006/01/02 15:04:05"),
			"updated_at": v.UpdatedAt,
			"author":     v.User.Username,
			"comment":    v.Reviews,
			"image_url":  imgUrl,
			"desc":       v.Description,
			"tags":       v.Tags,
			"id":         strconv.Itoa(int(v.ID)),
		})
	}

	// todo 取各分类的文章

	postSvc := post.NewService(c.logger, c.config, c.repository)
	res, _ := postSvc.Popular(ctx)

	total, _ := c.repository.Post().Count()

	// todo: 可以考虑从CACHE里取
	categories, _ := c.repository.Category().FindAll()

	// todo: 取各个分类下的头几篇文章
	for k, v := range categories {
		posts, _ := c.repository.Post().FindByCategoryId(v.Id, 7)
		categories[k].Posts = posts
	}

	categoryPosts := make(map[int64][]map[string]string)

	for _, v := range categories {
		for _, vv := range v.Posts {
			var imgUrl string
			if len(vv.Images) > 0 {
				imgUrl = c.config.GetString("server", "image_domain") + "/" + vv.Images[0].ImagePath
			}
			categoryPosts[v.Id] = append(categoryPosts[v.Id], map[string]string{
				"title":     vv.Title,
				"image_url": imgUrl,
				"desc":      vv.Description,
				"id":        strconv.Itoa(int(vv.ID)),
			})
		}
	}

	return map[string]interface{}{
		"stars":         starsData,
		"list":          posts,
		"populars":      res,
		"total":         total,
		"categories":    categories,
		"categoryPosts": categoryPosts,
	}, nil
}
