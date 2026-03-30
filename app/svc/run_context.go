package svc

import (
	"context"
)

// RunContext is an embeddable struct that implements RunContextAware.
// Services can embed this to easily gain access to the run context.
type RunContext struct {
	ctx context.Context
}

func (rc *RunContext) SetRunContext(ctx context.Context) {
	rc.ctx = ctx
}

// Context returns the run context.
// It returns nil if SetRunContext has not been called yet.
func (rc *RunContext) Context() context.Context {
	if rc.ctx == nil {
		return context.Background()
	}
	return rc.ctx
}
