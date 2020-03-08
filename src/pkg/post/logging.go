package post

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/shalom/src/repository"
	"github.com/icowan/shalom/src/repository/types"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
}

func (s *loggingService) AdminList(ctx context.Context, order, by, category, tag string, pageSize, offset int, keyword string) (posts []*types.Post, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "AdminList",
			"order", order,
			"by", by,
			"category", category,
			"tag", tag,
			"pageSize", pageSize,
			"offset", offset,
			"keyword", keyword,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AdminList(ctx, order, by, category, tag, pageSize, offset, keyword)
}

func (s *loggingService) Search(ctx context.Context, keyword, tag string, categoryId int64, offset, pageSize int) (posts []*types.Post, total int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Search",
			"keyword", keyword,
			"tag", tag,
			"categoryId", categoryId,
			"offset", offset,
			"pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Search(ctx, keyword, tag, categoryId, offset, pageSize)
}

func (s *loggingService) Get(ctx context.Context, id int64) (rs map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Get",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, id)
}

func (s *loggingService) Awesome(ctx context.Context, id int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Awesome",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Awesome(ctx, id)
}

func (s *loggingService) List(ctx context.Context, order, by, category string, pageSize, offset int) (rs []map[string]interface{}, count int64, other map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "List",
			"order", order,
			"by", by,
			"category", category,
			"pageSize", pageSize,
			"offset", offset,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, order, by, category, pageSize, offset)
}

func (s *loggingService) Popular(ctx context.Context) (rs []map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Popular",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Popular(ctx)
}

func (s *loggingService) NewPost(ctx context.Context, title, description, content string,
	postStatus repository.PostStatus, categoryNames, tagNames []string, markdown bool, imageId int64) (id int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "NewPost",
			"title", title,
			"description", description,
			"postStatus", postStatus,
			"categoryNames", categoryNames,
			"tagNames", tagNames,
			"markdown", markdown,
			"imageId", imageId,
			"postId", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.NewPost(ctx, title, description, content, postStatus, categoryNames, tagNames, markdown, imageId)
}

func (s *loggingService) Put(ctx context.Context, id int64, title, description, content string,
	postStatus repository.PostStatus, categoryNames, tagNames []string, markdown bool, imageId int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Put",
			"id", id,
			"title", title,
			"description", description,
			"postStatus", postStatus,
			"categoryNames", categoryNames,
			"tagNames", tagNames,
			"markdown", markdown,
			"imageId", imageId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Put(ctx, id, title, description, content, postStatus, categoryNames, tagNames, markdown, imageId)
}

func (s *loggingService) Delete(ctx context.Context, id int64) (err error) {
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

func (s *loggingService) Restore(ctx context.Context, id int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Restore",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Restore(ctx, id)
}

func (s *loggingService) Detail(ctx context.Context, id int64) (rs *types.Post, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Detail",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx, id)
}

func (s *loggingService) Star(ctx context.Context, id int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Star",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Star(ctx, id)
}
