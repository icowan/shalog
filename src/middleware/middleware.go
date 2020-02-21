package middleware

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	sjwt "github.com/icowan/blog/src/jwt"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

type ADSContext string

const (
	ContextUserEmail ADSContext = "user-email"
	ContextUsername  ADSContext = "user-name"
	ContextUserId    ADSContext = "user-id"
)

var (
	ErrorASD      = errors.New("权限验证失败！")
	ErrorNotLogin = errors.New("请先登录！")
)

func LoginMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			reqAuth := ctx.Value(kithttp.ContextKeyRequestAuthorization)
			if reqAuth == nil {
				return nil, ErrorNotLogin
			}

			token := reqAuth.(string)

			var claims sjwt.ArithmeticCustomClaims
			tk, err := jwt.ParseWithClaims(token, &claims, sjwt.JwtKeyFunc)
			if err != nil {
				_ = level.Error(logger).Log("jwt", "ParseWithClaims", "err", err.Error())
				return nil, errors.Wrap(err, ErrorASD.Error())
			}

			claim, ok := tk.Claims.(*sjwt.ArithmeticCustomClaims)
			if !ok {
				_ = level.Error(logger).Log("jwt", "Claims", "err", "")
				return
			}

			ctx = context.WithValue(ctx, ContextUserId, claim.UserId)
			ctx = context.WithValue(ctx, ContextUsername, claim.Username)
			ctx = context.WithValue(ctx, ContextUserEmail, claim.Email)

			return next(ctx, request)
		}
	}
}

var ErrLimitExceed = errors.New("Rate limit exceed!")

func TokenBucketLimitter(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
