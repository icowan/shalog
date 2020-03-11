package link

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/shalog/src/repository"
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

func (l *loggingServer) ApplyLink(ctx context.Context, name, link, icon string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "ApplyLink",
			"name", name,
			"link", link,
			"icon", icon,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.ApplyLink(ctx, name, link, icon)
}

func (l *loggingServer) Post(ctx context.Context, name, link, icon string, state repository.LinkState) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "Post",
			"name", name,
			"link", link,
			"icon", icon,
			"state", state,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.Post(ctx, name, link, icon, state)
}

func (l *loggingServer) Delete(ctx context.Context, id int64) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "Delete",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.Delete(ctx, id)
}

func (l *loggingServer) List(ctx context.Context) (links []*types.Link, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "List",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.List(ctx)
}

func (l *loggingServer) All(ctx context.Context) (links []*types.Link, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "All",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.All(ctx)
}

func (l *loggingServer) Pass(ctx context.Context, id int64) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "Pass",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.Pass(ctx, id)
}
