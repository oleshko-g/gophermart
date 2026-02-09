// Package oglog is a helper logging package
package oglog

import (
	"context"

	"goa.design/clue/log"
)

// NewLoggingCtx returns the context with goa logger
func NewLoggingCtx() context.Context {
	ctx := context.Background()

	opts := []log.LogOption{
		log.WithFormat(log.FormatTerminal),
		log.WithDebug(),
	}

	// puts logger in ctx
	logCtx := log.Context(ctx, opts...)
	return logCtx
}
