package gomock

import (
	"context"
	"net/http"

	"go.uber.org/mock/gomock"
	"gopkg.in/h2non/gock.v1"
)

type ctxMatcher struct {
	validate func(ctx context.Context) error
}

func ContextMatcher(validate func(ctx context.Context) error) gomock.Matcher {
	return &ctxMatcher{
		validate: validate,
	}
}

func (m *ctxMatcher) Matches(x interface{}) bool {
	ctx, ok := x.(context.Context)
	if !ok {
		return false
	}

	err := m.validate(ctx)
	if err != nil {
		return false
	}

	return err == nil
}

func (m *ctxMatcher) String() string {
	return "context matches"
}

// gockMatcher is a gock matcher for http.Request
type gockMatcher struct {
	gock gock.Mock
}

func GockMatcher(req *gock.Request) gomock.Matcher {
	if req.Response == nil {
		req.Response = gock.NewResponse()
	}

	return &gockMatcher{
		gock: gock.NewMock(req, req.Response),
	}
}

// Matches returns whether x is a match.
func (m *gockMatcher) Matches(x any) bool {
	req, ok := x.(*http.Request)
	if !ok {
		return false
	}

	match, err := m.gock.Match(req)
	if err != nil {
		return false
	}

	return match
}

// String describes what the matcher matches.
func (m *gockMatcher) String() string {
	return "HTTP request matching right method"
}
