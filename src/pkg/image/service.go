package image

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/blog/src/config"
	"github.com/icowan/blog/src/repository"
	"github.com/icowan/blog/src/repository/types"
)

type Service interface {
	List(ctx context.Context, pageSize, offset int) (images []types.Image, count int64, err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
	config     *config.Config
}

func (s *service) List(ctx context.Context, pageSize, offset int) (images []types.Image, count int64, err error) {
	images, count, err = s.repository.Image().FindAll(pageSize, offset)
	if err != nil {
		_ = level.Error(s.logger).Log("repository.Image", "FindAll", "err", err.Error())
		return
	}

	for k, v := range images {
		images[k].ImagePath = imageUrl(v.ImagePath, s.config.GetString("server", "image_domain"))
	}

	return
}

func NewService(logger log.Logger, repository repository.Repository, config *config.Config) Service {
	return &service{
		logger:     logger,
		repository: repository,
		config:     config,
	}
}

func imageUrl(path, imageDomain string) string {
	return imageDomain + "/" + path
}
