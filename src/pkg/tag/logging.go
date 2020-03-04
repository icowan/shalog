/**
 * @Time: 2020/3/1 16:01
 * @Author: solacowa@gmail.com
 * @File: middleware
 * @Software: GoLand
 */

package tag

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/blog/src/repository/types"
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

func (s *loggingServer) All(ctx context.Context) (categories []*types.Tag, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "All",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.All(ctx)
}

func (s *loggingServer) Delete(ctx context.Context, id int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Delete",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(ctx, id)
}

func (s *loggingServer) Put(ctx context.Context, id int64, name string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Put",
			"id", id,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Put(ctx, id, name)
}

func (s *loggingServer) Post(ctx context.Context, name string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Post",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, name)
}

func (s *loggingServer) Get(ctx context.Context, name string) (tag types.Tag, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Get",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, name)
}
