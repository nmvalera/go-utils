package tag

import (
	"context"
	"fmt"
	"strings"
)

func ExpectTagsOnContext(ctx context.Context, expectedTags ...*Tag) error {
	tags := FromContext(ctx)
	if len(tags) != len(expectedTags) {
		return fmt.Errorf("expected %d tags, got %d", len(expectedTags), len(tags))
	}

	tagErrors := make([]string, 0)
	for i, tag := range tags {
		if tag.Key != expectedTags[i].Key || tag.Value.String() != expectedTags[i].Value.String() {
			tagErrors = append(tagErrors, fmt.Sprintf("tag %v does not match: %v", tag, expectedTags[i]))
		}
	}

	if len(tagErrors) > 0 {
		return fmt.Errorf("tags do not match: %s", strings.Join(tagErrors, "\n"))
	}

	return nil
}
