package svc

import (
	"context"
	"testing"

	"github.com/nmvalera/go-utils/tag"
	"github.com/stretchr/testify/require"
)

func TestTagged(t *testing.T) {
	type service struct {
		*Tagged
	}

	testSvc := &service{Tagged: NewTagged()}

	testSvc.WithTags(tag.Key("test-key").String("test-value"))

	ctx := testSvc.Context(context.Background())
	tags := tag.FromContext(ctx)
	require.Equal(t, 1, len(tags))
	require.Equal(t, "test-key", string(tags[0].Key))
	require.Equal(t, "test-value", tags[0].Value.String())
}
