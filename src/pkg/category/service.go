/**
 * @Time: 2020/2/29 17:10
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package category

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/icowan/blog/src/repository"
	"github.com/icowan/blog/src/repository/types"
	"github.com/pkg/errors"
)

var (
	ErrCategoryParams    = errors.New("参数错误!")
	ErrCategoryParamName = errors.New("名称不能为空")
)

type Service interface {
	List(ctx context.Context) (categories []*types.Category, err error)
	Post(ctx context.Context, title, description string, parentId int64) (err error)
	Delete(ctx context.Context, id int64) (err error)
	Put(ctx context.Context, id int64, title, description string, parentId int64) (err error)
}

type service struct {
	repository repository.Repository
	logger     log.Logger
}

func (s *service) List(ctx context.Context) (categories []*types.Category, err error) {
	categories, err = s.repository.Category().FindAll()
	if err != nil {
		return
	}

	for k, v := range categories {
		categories[k].PostCount = s.repository.Category().CountPosts(v)
	}

	return
}

func (s *service) Post(ctx context.Context, title, description string, parentId int64) (err error) {
	return s.repository.Category().Add(title, description, parentId)
}

func (s *service) Delete(ctx context.Context, id int64) (err error) {
	return s.repository.Category().Delete(id)
}

func (s *service) Put(ctx context.Context, id int64, title, description string, parentId int64) (err error) {
	return s.repository.Category().Put(id, title, description, parentId)
}

func NewService(logger log.Logger, repository repository.Repository) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}
