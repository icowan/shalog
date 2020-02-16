package post

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/nsini/blog/src/repository"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
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

var ErrLimitExceed = errors.New("Rate limit exceed!")

func newTokenBucketLimitter(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
