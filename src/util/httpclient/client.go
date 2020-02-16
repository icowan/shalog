/**
 * @Time : 2019-09-06 10:43
 * @Author : solacowa@gmail.com
 * @File : client
 * @Software: GoLand
 */

package httpclient

import (
	"context"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func NewClient(ctx context.Context, instance, method string, req interface{}) (body interface{}, err error) {

	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}

	tgt, err := url.Parse(instance)
	if err != nil {
		return
	}

	ep := kithttp.NewClient(method, tgt, func(ctx context.Context, r *http.Request, i interface{}) error {
		fmt.Println(i)
		fmt.Println(r)
		fmt.Println(ctx)
		fmt.Println("enc-------------")
		return nil
	}, func(ctx context.Context, res *http.Response) (response interface{}, err error) {
		b, _ := ioutil.ReadAll(res.Body)
		return string(b), nil
	}, kithttp.ClientAfter(func(ctx context.Context, res *http.Response) context.Context {
		fmt.Println(res.StatusCode)
		return ctx
	}), kithttp.ClientBefore(func(ctx context.Context, r *http.Request) context.Context {
		return ctx
	})).Endpoint()

	res, err := ep(ctx, req)
	if err != nil {
		return
	}

	return res, nil
}
