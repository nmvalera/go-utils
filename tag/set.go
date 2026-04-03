package tag

import "fmt"

var EmptySet = Set{}

// Set represents a immutable set of
type Set []*Tag

// WithTags returns a new set with the given tags added to the set. If a tag with the same key already exists in the set,
// the new tag will replace the old tag.
func (s Set) WithTags(tags ...*Tag) Set {
	result := make(Set, len(s))
	for i, t := range s {
		result[i] = t.Copy()
	}

	newTags := make(Set, 0)
	for _, tag := range tags {
		var existed bool
		for i, oldTag := range result {
			if oldTag.Key == tag.Key {
				result[i] = merge(oldTag, tag)
				existed = true
				break
			}
		}
		if !existed {
			newTags = append(newTags, tag.Copy())
		}
	}

	return append(result, newTags...)
}

func merge(oldTag, newTag *Tag) *Tag {
	if newTag.chained != nil && *newTag.chained {
		if oldTag.Value.Type == STRING {
			return &Tag{
				Key:   oldTag.Key,
				Value: StringValue(oldTag.Value.Interface.(string) + "." + newTag.Value.Interface.(string)),
			}
		}

		if oldTag.Value.Type == MAP && newTag.Value.Type == MAP {
			oldTags := oldTag.Value.Interface.(Set)
			newTags := newTag.Value.Interface.(Set)
			return &Tag{
				Key:   oldTag.Key,
				Value: MapValue(oldTags.WithTags(newTags...)...),
			}
		}

		if oldTag.Value.Type == MAP {
			return newTag.Copy()
		}

		panic(fmt.Sprintf("cannot chain %T with %T for key %s", oldTag.Value.Type, newTag.Value.Type, oldTag.Key))
	}
	return newTag.Copy()
}
