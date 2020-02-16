package api

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Post(ctx context.Context, req postRequest) (rs newPostResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Post",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, req)
}

func (s *loggingService) UploadImage(ctx context.Context, req uploadImageRequest) (rs imageResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "UploadImage",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.UploadImage(ctx, req)
}
