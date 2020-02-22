package link

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/icowan/blog/src/repository"
	"github.com/icowan/blog/src/repository/types"
	"github.com/pkg/errors"
)

var (
	ErrParams = errors.New("参数错误：")
)

type Service interface {
	// 申请友链
	ApplyLink(ctx context.Context, name, link, icon string) (err error)

	// 增加友链
	Post(ctx context.Context, name, link, icon string, state int) (err error)

	// 删除
	Delete(ctx context.Context, id int64) (err error)

	// 列表
	List(ctx context.Context) (links []*types.Link, err error)
}

type service struct {
	repository repository.Repository
	logger     log.Logger
}

func (s *service) ApplyLink(ctx context.Context, name, link, icon string) (err error) {
	panic("implement me")
}

func (s *service) Post(ctx context.Context, name, link, icon string, state int) (err error) {
	panic("implement me")
}

func (s *service) Delete(ctx context.Context, id int64) (err error) {
	panic("implement me")
}

func (s *service) List(ctx context.Context) (links []*types.Link, err error) {
	panic("implement me")
}

func NewService(logger log.Logger, repository repository.Repository) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}
