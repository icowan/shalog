package link

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/blog/src/repository"
	"github.com/icowan/blog/src/repository/types"
	"github.com/pkg/errors"
)

var (
	ErrParams       = errors.New("参数错误")
	ErrParamLinkLen = errors.New("参数错误,链接太长")
	ErrParamNameLen = errors.New("参数错误,名称太长")
	ErrLinkNotfound = errors.New("未查到数据")
	ErrLinkUpdate   = errors.New("更新错误")
	ErrLinkAdd      = errors.New("添加失败")
)

type Service interface {
	// 申请友链
	ApplyLink(ctx context.Context, name, link, icon string) (err error)

	// 增加友链
	Post(ctx context.Context, name, link, icon string, state repository.LinkState) (err error)

	// 删除
	Delete(ctx context.Context, id int64) (err error)

	// 列表
	List(ctx context.Context) (links []*types.Link, err error)

	// 通过
	Pass(ctx context.Context, id int64) (err error)
}

type service struct {
	repository repository.Repository
	logger     log.Logger
}

func (s *service) Pass(ctx context.Context, id int64) (err error) {
	link, err := s.repository.Link().Find(id)
	if err != nil {
		_ = level.Error(s.logger).Log("repository.Link", "Find", "id", id, "err", err.Error())
		err = ErrLinkNotfound
		return
	}
	link.State = 1
	if err = s.repository.Link().Update(&link); err != nil {
		err = errors.Wrap(err, ErrLinkUpdate.Error())
	}
	return
}

func (s *service) ApplyLink(ctx context.Context, name, link, icon string) (err error) {
	// 在 transport的request 处理link链接

	err = s.repository.Link().Add(name, link, icon, repository.LinkStateApply)
	if err != nil {
		_ = level.Error(s.logger).Log("repository.Link", "Add", "err", err.Error())
		err = errors.Wrap(err, ErrLinkAdd.Error())
	}

	return
}

func (s *service) Post(ctx context.Context, name, link, icon string, state repository.LinkState) (err error) {
	err = s.repository.Link().Add(name, link, icon, state)
	if err != nil {
		_ = level.Error(s.logger).Log("repository.Link", "Add", "err", err.Error())
		err = errors.Wrap(err, ErrLinkAdd.Error())
	}
	return
}

func (s *service) Delete(ctx context.Context, id int64) (err error) {
	return s.repository.Link().Delete(id)
}

func (s *service) List(ctx context.Context) (links []*types.Link, err error) {
	return s.repository.Link().List()
}

func NewService(logger log.Logger, repository repository.Repository) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}
