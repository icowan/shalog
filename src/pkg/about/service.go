package about

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Service interface {
	About(ctx context.Context) (rs map[string]interface{}, err error)
}

type service struct {
	logger log.Logger
}

func NewService(logger log.Logger) Service {
	return &service{
		logger: logger,
	}
}

func (c *service) About(ctx context.Context) (rs map[string]interface{}, err error) {

	return
}
