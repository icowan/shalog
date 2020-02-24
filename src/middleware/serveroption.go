package middleware

import (
	"context"
	"net/http"
)

func SettingsRequest(settings map[string]string) func(context.Context, *http.Request) context.Context {
	return func(ctx context.Context, r *http.Request) context.Context {
		//for key, val := range settings {
		//	ctx = context.WithValue(ctx, key, val)
		//}
		ctx = context.WithValue(ctx, "settings", settings)
		return ctx
	}
}
