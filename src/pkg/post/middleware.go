package post

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/icowan/shalom/src/repository"
)

func UpdatePostReadNum(logger log.Logger, repository repository.Repository) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			go func() {
				req := request.(postRequest)
				detail, err := repository.Post().FindOnce(req.Id)
				if err != nil {
					return
				}
				if err = repository.Post().SetReadNum(detail); err != nil {
					_ = logger.Log("post.SetReadNum", err.Error())
				}
			}()

			// read cache

			return next(ctx, request)
		}
	}
}
