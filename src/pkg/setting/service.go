package setting

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/icowan/blog/src/config"
	"github.com/icowan/blog/src/repository"
	"github.com/icowan/blog/src/repository/types"
	"github.com/pkg/errors"
)

var (
	ErrSettingNotfound  = errors.New("配置的关键字错误")
	ErrSettingReqParams = errors.New("参数错误")
)

type Service interface {
	// 获取配置信息
	Get(ctx context.Context, key repository.SettingKey) (setting types.Setting, err error)

	// 创建配置
	Post(ctx context.Context, key repository.SettingKey, val, description string) (err error)

	// 删除配置
	Delete(ctx context.Context, key repository.SettingKey) (err error)

	// 更新配置
	Put(ctx context.Context, key repository.SettingKey, val, description string) (err error)

	// 上传图片
	UploadImage(ctx context.Context, key repository.SettingKey, imgUrl string, file *imageFile) (urlAddr string, err error)

	// 配置列表
	List(ctx context.Context) (settings []*types.Setting, err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
	config     *config.Config
}

func (s *service) Get(ctx context.Context, key repository.SettingKey) (setting types.Setting, err error) {
	return s.repository.Setting().Find(key)
}

func (s *service) Post(ctx context.Context, key repository.SettingKey, val, description string) (err error) {
	return s.repository.Setting().Add(key, val, description)
}

func (s *service) Delete(ctx context.Context, key repository.SettingKey) (err error) {
	return s.repository.Setting().Delete(key)
}

func (s *service) Put(ctx context.Context, key repository.SettingKey, val, description string) (err error) {
	setting, err := s.repository.Setting().Find(key)
	if err != nil {
		err = errors.Wrap(err, ErrSettingNotfound.Error())
		return
	}
	setting.Description = description
	setting.Value = val
	return s.repository.Setting().Update(&setting)
}

func (s *service) UploadImage(ctx context.Context, key repository.SettingKey, imgUrl string, file *imageFile) (urlAddr string, err error) {
	panic("implement me")
}

func (s *service) List(ctx context.Context) (settings []*types.Setting, err error) {
	return s.repository.Setting().List()
}

func NewService(logger log.Logger, repository repository.Repository, config *config.Config) Service {
	return &service{
		logger:     logger,
		repository: repository,
		config:     config,
	}
}
