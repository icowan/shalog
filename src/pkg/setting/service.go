package setting

import (
	"context"
	"github.com/chanxuehong/wechat/mp/core"
	"github.com/chanxuehong/wechat/mp/menu"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/config"
	"github.com/icowan/shalog/src/repository"
	"github.com/icowan/shalog/src/repository/types"
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

	// 微信菜单
	WechatMenu(ctx context.Context) (err error)
}

type service struct {
	logger       log.Logger
	repository   repository.Repository
	config       *config.Config
	wechatClient *core.Client
	traceId      string
}

func (s *service) WechatMenu(ctx context.Context) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "WechatMenu")

	buttons := []menu.Button{
		{
			Name:     "薛定谔的猿",
			AppId:    "wx7a80548009a2d40a",
			PagePath: "/pages/dashboard/dashboard",
		},
		{
			Name: "有钱得",
			//AppId:      "wx00dfa6eafd4162ac",
			//PagePath:   "/pages/index/index?scene=37e9a7ba7d018b8fd39047afe1469ecc",
			SubButtons: []menu.Button{
				{
					Name:     "饿了么返现",
					AppId:    "wx00dfa6eafd4162ac",
					PagePath: "/pages/index/index?scene=37e9a7ba7d018b8fd39047afe1469ecc",
				},
			},
		},
	}

	menuId, err := menu.AddConditionalMenu(s.wechatClient, &menu.Menu{
		Buttons: buttons,
	})
	if err != nil {
		_ = level.Error(logger).Log("menu", "AddConditionalMenu", "err", err.Error())
		return
	}

	_ = level.Debug(logger).Log("menuId", menuId)

	return
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
	if description != "" {
		setting.Description = description
	}
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
	var (
		accessTokenServer core.AccessTokenServer = core.NewDefaultAccessTokenServer(
			config.GetString("wechat", "app_id"),
			config.GetString("wechat", "app_secret"), nil)
		wechatClient *core.Client = core.NewClient(accessTokenServer, nil)
	)

	return &service{
		logger:       logger,
		repository:   repository,
		config:       config,
		wechatClient: wechatClient,
	}
}
