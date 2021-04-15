package setting

import (
	"context"
	"encoding/json"
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
	var buttons []menu.Button
	if st, err := s.repository.Setting().Find(repository.SettingWechatOfficialMenu); err == nil {
		if err = json.Unmarshal([]byte(st.Value), &buttons); err != nil {
			_ = level.Error(logger).Log("json", "Unmarshal", "err", err.Error())
			return err
		}
	} else {
		return err
	}

	if buttons == nil {
		return
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
	//var proxy func(r *http.Request) (*url.URL, error)
	//proxy = func(_ *http.Request) (*url.URL, error) {
	//	return url.Parse("http://127.0.0.1:1087")
	//}
	//
	//dialer := &net.Dialer{
	//	Timeout:   time.Duration(5 * int64(time.Second)),
	//	KeepAlive: time.Duration(5 * int64(time.Second)),
	//}
	//
	//cli := &http.Client{
	//	Transport: &http.Transport{
	//		Proxy: proxy, DialContext: dialer.DialContext,
	//		TLSClientConfig: &tls.Config{
	//			InsecureSkipVerify: false,
	//		},
	//	},
	//}

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
