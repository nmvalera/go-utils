package jsonrpc

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/kkrt-labs/go-utils/log"
	"go.uber.org/zap"
)

// ClientDecorator is a function that enable to decorate a JSON-RPC client with additional functionality
type ClientDecorator func(Client) Client

// WithVersion automatically set JSON-RPC request version
func WithVersion(v string) ClientDecorator {
	return func(c Client) Client {
		return ClientFunc(func(ctx context.Context, req *Request, res any) error {
			req.Version = v
			return c.Call(ctx, req, res)
		})
	}
}

// WithIncrementalID automatically increments JSON-RPC request ID
func WithIncrementalID() ClientDecorator {
	var idCounter uint32
	return func(c Client) Client {
		return ClientFunc(func(ctx context.Context, req *Request, res any) error {
			req.ID = atomic.AddUint32(&idCounter, 1) - 1
			return c.Call(ctx, req, res)
		})
	}
}

// WithExponentialBackOffRetry automatically retries JSON-RPC calls
func WithExponentialBackOffRetry(opts ...backoff.ExponentialBackOffOpts) ClientDecorator {
	pool := &sync.Pool{
		New: func() any {
			return backoff.NewExponentialBackOff(opts...)
		},
	}
	return func(c Client) Client {
		return ClientFunc(func(ctx context.Context, req *Request, res any) error {
			bckff := pool.Get().(*backoff.ExponentialBackOff)
			defer func() {
				bckff.Reset()
				pool.Put(bckff)
			}()

			attempt := 0
			attemptReq := req
			return backoff.RetryNotify(
				func() error {
					return c.Call(ctx, attemptReq, res)
				},
				backoff.WithContext(bckff, ctx),
				func(err error, d time.Duration) {
					attempt++
					// We need to increment the ID for each retry attempt
					// so that we don't possibly overwrite the response of the previous attempt
					attemptReq = &Request{
						Method:  req.Method,
						Version: req.Version,
						Params:  req.Params,
						ID:      fmt.Sprintf("%s#%d", req.ID, attempt),
					}
					log.LoggerFromContext(ctx).Warn(
						fmt.Sprintf("Call failed, retrying in %s...", d),
						zap.Error(err),
					)
				},
			)
		})
	}
}

// WithTimeout automatically sets a timeout for JSON-RPC calls
func WithTimeout(d time.Duration) ClientDecorator {
	return func(c Client) Client {
		return ClientFunc(func(ctx context.Context, req *Request, res any) error {
			deadline := time.Now().Add(d)

			cancelCtx, cancel := context.WithDeadline(ctx, deadline)
			defer cancel()

			err := c.Call(cancelCtx, req, res)
			if err != nil && time.Now().After(deadline) {
				err = fmt.Errorf("jsonrpc: call timed out after %q: %w", d, err)
			}
			return err
		})
	}
}
