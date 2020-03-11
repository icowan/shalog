package account

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/shalog/src/repository/types"
	"time"
)

type loggingServer struct {
	logger log.Logger
	Service
}

func NewLoggingServer(logger log.Logger, s Service) Service {
	return &loggingServer{
		logger:  level.Info(logger),
		Service: s,
	}
}

func (s *loggingServer) Get(ctx context.Context) (user types.User, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Get",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx)
}

func (s *loggingServer) Logout(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Logout",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Logout(ctx)
}

func (s *loggingServer) Put(ctx context.Context, username, usernameCanonical, email, emailCanonical, password string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Put",
			"username", username,
			"email", email,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Put(ctx, username, usernameCanonical, email, emailCanonical, password)
}

func (s *loggingServer) Login(ctx context.Context, username, password string, checkbox bool) (user string, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Login",
			"username", username,
			"checkbox", checkbox,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Login(ctx, username, password, checkbox)
}
