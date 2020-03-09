package setting

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/shalom/src/repository"
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

func (l *loggingServer) List(ctx context.Context) (settings []*types.Setting, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "ApplyLink",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.List(ctx)
}

func (l *loggingServer) Get(ctx context.Context, key repository.SettingKey) (settings types.Setting, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "Get",
			"key", key,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.Get(ctx, key)
}

func (l *loggingServer) Delete(ctx context.Context, key repository.SettingKey) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "Delete",
			"key", key,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.Delete(ctx, key)
}

func (l *loggingServer) Put(ctx context.Context, key repository.SettingKey, val, description string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "Put",
			"key", key,
			"val", val,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.Put(ctx, key, val, description)
}

func (l *loggingServer) Post(ctx context.Context, key repository.SettingKey, val, description string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "Post",
			"key", key,
			"val", val,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.Post(ctx, key, val, description)
}

func (l *loggingServer) UploadImage(ctx context.Context, key repository.SettingKey, imgUrl string, file *imageFile) (urlAddr string, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			"method", "UploadImage",
			"key", key,
			"imgUrl", imgUrl,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.Service.UploadImage(ctx, key, imgUrl, file)
}
