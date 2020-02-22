package post

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/blog/src/repository"
	"github.com/icowan/blog/src/repository/types"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{level.Info(logger), s}
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
	postStatus repository.PostStatus, categoryIds, tagIds []int64, markdown bool, imageId int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "NewPost",
			"title", title,
			"description", description,
			"postStatus", postStatus,
			"categoryIds", categoryIds,
			"tagIds", tagIds,
			"markdown", markdown,
			"imageId", imageId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.NewPost(ctx, title, description, content, postStatus, categoryIds, tagIds, markdown, imageId)
}

func (s *loggingService) Put(ctx context.Context, id int64, title, description, content string,
	postStatus repository.PostStatus, categoryIds, tagIds []int64, markdown bool, imageId int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Put",
			"id", id,
			"title", title,
			"description", description,
			"postStatus", postStatus,
			"categoryIds", categoryIds,
			"tagIds", tagIds,
			"markdown", markdown,
			"imageId", imageId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Put(ctx, id, title, description, content, postStatus, categoryIds, tagIds, markdown, imageId)
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
