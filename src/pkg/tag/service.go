/**
 * @Time: 2020/3/1 15:51
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package tag

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/icowan/shalom/src/repository"
	"github.com/icowan/shalom/src/repository/types"
	"github.com/pkg/errors"
)

var (
	ErrTagParams = errors.New("参数错误")
)

type Service interface {
	All(ctx context.Context) (tags []*types.Tag, err error)
	Post(ctx context.Context, name string) (err error)
	Put(ctx context.Context, id int64, name string) (err error)
	Delete(ctx context.Context, id int64) (err error)
	Get(ctx context.Context, name string) (tags types.Tag, err error)
	List(ctx context.Context, tagName string, limit, offset int) (tags []*types.Tag, count int64, err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func (s *service) List(ctx context.Context, tagName string, limit, offset int) (tags []*types.Tag, count int64, err error) {
	tags, count, err = s.repository.Tag().List(tagName, limit, offset)
	return
}

func (s *service) All(ctx context.Context) (tags []*types.Tag, err error) {
	return s.repository.Tag().All(50)
}

func (s *service) Post(ctx context.Context, name string) (err error) {
	_, err = s.repository.Tag().FirstOrCreate(name)
	return
}

func (s *service) Put(ctx context.Context, id int64, name string) (err error) {
	return s.repository.Tag().Update(id, name)
}

func (s *service) Delete(ctx context.Context, id int64) (err error) {
	return s.repository.Tag().Delete(id)
}

func (s *service) Get(ctx context.Context, name string) (tags types.Tag, err error) {
	panic("implement me")
}

func NewService(logger log.Logger, repository repository.Repository) Service {
	return &service{
		logger:     logger,
		repository: repository,
	}
}
