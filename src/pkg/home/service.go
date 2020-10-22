package home

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/icowan/shalog/src/pkg/post"
	"github.com/icowan/shalog/src/repository"
	"github.com/icowan/shalog/src/repository/types"
	"github.com/pkg/errors"
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

type indexCh struct {
	List  []map[string]interface{}
	Error error
}

func (c *service) getStars(ch chan indexCh) {
	var res []map[string]interface{}
	stars, err := c.repository.Post().Stars()
	if err != nil {
		err = errors.Wrap(err, "Post Stars")
		ch <- indexCh{
			List:  nil,
			Error: err,
		}
		return
	}
	for _, v := range stars {
		var imgUrl string
		if len(v.Images) > 0 {
			imgUrl = c.config.GetString(config.SectionServer, repository.SettingGlobalDomainImage.String()) + "/" + v.Images[0].ImagePath
		}
		res = append(res, map[string]interface{}{
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
	ch <- indexCh{
		List:  res,
		Error: nil,
	}
}
func (c *service) getPosts(ch chan indexCh) {
	var res []map[string]interface{}
	var err error
	if list, err := c.repository.Post().Index(); err == nil {
		for _, v := range list {
			var imgUrl string
			if len(v.Images) > 0 {
				imgUrl = c.config.GetString(config.SectionServer, repository.SettingGlobalDomainImage.String()) + "/" + v.Images[0].ImagePath
			}
			res = append(res, map[string]interface{}{
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
	} else {
		err = level.Error(c.logger).Log("Post", "Index", "err", err.Error())
	}
	ch <- indexCh{
		List:  res,
		Error: err,
	}
}

func (c *service) Index(ctx context.Context) (rs map[string]interface{}, err error) {
	starsCh := make(chan indexCh)
	postsCh := make(chan indexCh)
	popularCh := make(chan indexCh)
	totalCh := make(chan int64)
	categoryPostsCh := make(chan map[int64][]map[string]string)
	go c.getStars(starsCh)
	go c.getPosts(postsCh)

	go func() {
		postSvc := post.NewService(c.logger, c.config, c.repository)
		res, err := postSvc.Popular(ctx)
		popularCh <- indexCh{
			List:  res,
			Error: err,
		}
	}()

	go func() {
		total, _ := c.repository.Post().Count()
		totalCh <- total
	}()

	var categories []*types.Category
	{
		go func() {
			categories, _ = c.repository.Category().FindAll()

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
						imgUrl = c.config.GetString(config.SectionServer, repository.SettingGlobalDomainImage.String()) + "/" + vv.Images[0].ImagePath
					}
					categoryPosts[v.Id] = append(categoryPosts[v.Id], map[string]string{
						"title":     vv.Title,
						"image_url": imgUrl,
						"desc":      vv.Description,
						"id":        strconv.Itoa(int(vv.ID)),
					})
				}
			}
			categoryPostsCh <- categoryPosts
		}()
	}

	stars := <-starsCh
	posts := <-postsCh
	popular := <-popularCh
	total := <-totalCh
	categoriesPosts := <-categoryPostsCh
	close(starsCh)
	close(postsCh)
	close(popularCh)
	close(totalCh)
	close(categoryPostsCh)

	if stars.Error != nil {
		_ = level.Error(c.logger).Log("err", stars.Error.Error())
	}
	if posts.Error != nil {
		_ = level.Error(c.logger).Log("err", posts.Error.Error())
	}

	return map[string]interface{}{
		"stars":         stars.List,
		"list":          posts.List,
		"populars":      popular.List,
		"total":         total,
		"categories":    categories,
		"categoryPosts": categoriesPosts,
	}, nil
}
