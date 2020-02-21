package account

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/blog/src/config"
	"github.com/icowan/blog/src/encode"
	sjwt "github.com/icowan/blog/src/jwt"
	"github.com/icowan/blog/src/middleware"
	"github.com/icowan/blog/src/repository"
	"github.com/icowan/blog/src/repository/types"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrBodyParams    = errors.New("参数错误!")
	ErrParamsNotNull = errors.New("用户名或密码不能为空!")
	ErrUserOrPwd     = errors.New("用户名或密码错误!")
	ErrContextId     = errors.New("用户不存在！")
	ErrUserUpdate    = errors.New("用户更新错误！")
)

type Service interface {
	// 登录
	Login(ctx context.Context, username, password string, checkbox bool) (res string, err error)

	// 退出
	Logout(ctx context.Context) (err error)

	// 获取用户信息
	Get(ctx context.Context) (user types.User, err error)

	// 更新用户信息
	Put(ctx context.Context, username, usernameCanonical, email, emailCanonical, password string) (err error)
}

type service struct {
	repository repository.Repository
	logger     log.Logger
	config     *config.Config
}

func (s *service) Get(ctx context.Context) (user types.User, err error) {
	id, ok := ctx.Value(middleware.ContextUserId).(int64)
	if !ok {
		_ = level.Error(s.logger).Log("ctx", "Value", "err", "ContextUserId is not ok!")
		err = errors.New("ContextUserId is not ok!")
		return
	}
	user, err = s.repository.User().FindById(id)
	if err != nil {
		_ = level.Error(s.logger).Log("repository.User", "FindById", "id", id, "err", err.Error())
		return
	}

	return
}

func (s *service) Put(ctx context.Context, username, usernameCanonical, email, emailCanonical, password string) (err error) {
	id, ok := ctx.Value(middleware.ContextUserId).(int64)
	if !ok {
		_ = level.Error(s.logger).Log("ctx", "Value", "err", "ContextUserId is not ok!")
		err = errors.New("ContextUserId is not ok!")
		return
	}

	user, err := s.repository.User().FindById(id)
	if err != nil {
		_ = level.Error(s.logger).Log("repository.User", "FindById", "id", id, "err", err.Error())
		return errors.Wrap(err, ErrContextId.Error())
	}

	user.Email = email
	user.EmailCanonical = email
	user.Username = username
	user.UsernameCanonical = usernameCanonical
	user.Password = encode.EncodePassword(password, s.config.GetString("server", "app_key"))
	user.PasswordRequestedAt = time.Now()

	err = s.repository.User().Update(&user)
	if err != nil {
		_ = level.Error(s.logger).Log("repository.User", "Update", "err", err.Error())
		err = errors.Wrap(err, ErrUserUpdate.Error())
	}
	return
}

func (s *service) Login(ctx context.Context, username, password string, checkbox bool) (res string, err error) {
	user, err := s.repository.User().FindAndPwd(username, encode.EncodePassword(password, s.config.GetString("server", "app_key")))
	if err != nil {
		_ = level.Error(s.logger).Log("repository.User", "FindAndPwd", "err", err.Error())
		return "", errors.Wrap(err, ErrUserOrPwd.Error())
	}

	sessionTimeout, err := s.config.Int64("server", "session_timeout")
	if err != nil {
		sessionTimeout = 3600
	}
	expAt := time.Now().Add(time.Duration(sessionTimeout) * time.Second).Unix()

	claim := sjwt.ArithmeticCustomClaims{
		UserId:   user.ID,
		Username: user.Username,
		Email:    user.EmailCanonical,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    "system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	res, err = token.SignedString([]byte(sjwt.GetJwtKey()))
	if err != nil {
		_ = level.Error(s.logger).Log("token", "SignedString", "sjwt.GetJwtKey()", sjwt.GetJwtKey(), "err", err.Error())
	}

	// todo: token 考虑入cache

	return
}

func (s *service) Logout(ctx context.Context) (err error) {
	id, ok := ctx.Value(middleware.ContextUserId).(int64)
	if !ok {
		_ = level.Error(s.logger).Log("ctx", "Value", "err", "ContextUserId is not ok!")
		err = errors.New("ContextUserId is not ok!")
		return
	}

	_ = level.Debug(s.logger).Log("id", id)
	// todo 清除cache
	return
}

func NewService(logger log.Logger, repository repository.Repository, config *config.Config) Service {

	return &service{
		repository: repository,
		logger:     logger,
		config:     config,
	}
}
