package svc

import (
	"context"
	"sync"

	"github.com/kkrt-labs/go-utils/tag"
)

type Tagged struct {
	sync.RWMutex
	tag.Set
}

func NewTagged(tags ...*tag.Tag) *Tagged {
	return &Tagged{Set: tag.Set(tags)}
}

func (t *Tagged) WithTags(tags ...*tag.Tag) {
	t.Lock()
	t.Set = t.Set.WithTags(tags...)
	t.Unlock()
}

func (t *Tagged) Context(ctx context.Context, tags ...*tag.Tag) context.Context {
	t.RLock()
	defer t.RUnlock()
	return tag.WithTags(ctx, append(t.Set, tags...)...)
}
