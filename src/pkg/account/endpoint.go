package account

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/blog/src/encode"
)

type Endpoints struct {
	GetEndpoint    endpoint.Endpoint
	PutEndpoint    endpoint.Endpoint
	LoginEndpoint  endpoint.Endpoint
	LogoutEndpoint endpoint.Endpoint
}

type (
	putRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Checkbox bool   `json:"checkbox"`
	}
)

func NewEndpoint(s Service, mdw map[string][]endpoint.Middleware) Endpoints {

	eps := Endpoints{
		GetEndpoint:    makeGetEndpoint(s),
		PutEndpoint:    makePutEndpoint(s),
		LoginEndpoint:  makeLoginEndpoint(s),
		LogoutEndpoint: makeLogoutEndpoint(s),
	}

	for _, m := range mdw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}
	for _, m := range mdw["Put"] {
		eps.PutEndpoint = m(eps.PutEndpoint)
	}
	for _, m := range mdw["Login"] {
		eps.LoginEndpoint = m(eps.LoginEndpoint)
	}
	for _, m := range mdw["Logout"] {
		eps.LogoutEndpoint = m(eps.LogoutEndpoint)
	}

	return eps
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, err := s.Get(ctx)
		return encode.Response{
			Data:  user,
			Error: err,
		}, err
	}
}

func makePutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putRequest)
		err = s.Put(ctx, req.Username, req.Username, req.Email, req.Email, req.Password)
		return encode.Response{Error: err}, err
	}
}

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		res, err := s.Login(ctx, req.Username, req.Password, req.Checkbox)
		return encode.Response{
			Data: map[string]interface{}{
				"username": req.Username,
				"token":    res,
			},
			Error: err,
		}, err
	}
}

func makeLogoutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		err = s.Logout(ctx)
		return encode.Response{Error: err}, err
	}
}
