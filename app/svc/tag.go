package svc

import (
	"context"

	"github.com/nmvalera/go-utils/tag"
)

// Tagged enable to attach tags to a service
type Tagged struct {
	set tag.Set
}

func (t *Tagged) WithTags(tags ...*tag.Tag) {
	if t.set == nil {
		t.set = tag.EmptySet
	}
	t.set = t.set.WithTags(tags...)
}

func (t *Tagged) Context(ctx context.Context, tags ...*tag.Tag) context.Context {
	return tag.WithTags(ctx, append(t.set, tags...)...)
}
