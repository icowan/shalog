/**
 * @Time: 2020/2/29 17:29
 * @Author: solacowa@gmail.com
 * @File: middleware
 * @Software: GoLand
 */

package category

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/shalom/src/repository/types"
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

func (s *loggingServer) List(ctx context.Context) (categories []*types.Category, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "List",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx)
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

func (s *loggingServer) Put(ctx context.Context, id int64, name, description string, parentId int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Put",
			"id", id,
			"name", name,
			"description", description,
			"parentId", parentId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Put(ctx, id, name, description, parentId)
}

func (s *loggingServer) Post(ctx context.Context, name, description string, parentId int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Post",
			"name", name,
			"description", description,
			"parentId", parentId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, name, description, parentId)
}
