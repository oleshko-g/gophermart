package service

import (
	"context"
	"fmt"

	"goa.design/clue/log"
	goa "goa.design/goa/v3/pkg"
)

func WithLogEndpoint(endpoint goa.Endpoint) goa.Endpoint {
	return func(ctx context.Context, payload any) (result any, err error) {
		log.MustContainLogger(ctx)

		p := fmt.Sprintf("%+v", payload)
		log.Debug(ctx, log.KV{K: "payload", V: p})

		result, err = endpoint(ctx, payload)

		r := fmt.Sprintf("%+v", result)
		e := fmt.Sprintf("%+v", err)
		log.Debug(ctx, log.KV{K: "result", V: r}, log.KV{K: "error", V: e})

		return result, err
	}
}
