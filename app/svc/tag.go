package svc

import (
	"context"
	"sync"

	"github.com/nmvalera/go-utils/tag"
)

type Tagged struct {
	mux sync.RWMutex
	tag.Set
}

func NewTagged(tags ...*tag.Tag) *Tagged {
	return &Tagged{Set: tag.Set(tags)}
}

func (t *Tagged) WithTags(tags ...*tag.Tag) {
	t.mux.Lock()
	t.Set = t.Set.WithTags(tags...)
	t.mux.Unlock()
}

func (t *Tagged) Context(ctx context.Context, tags ...*tag.Tag) context.Context {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return tag.WithTags(ctx, append(t.Set, tags...)...)
}
