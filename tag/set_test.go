package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithTags(t *testing.T) {
	set := EmptySet.WithTags(
		Key("key1").String("value1.0"),
		Key("key2").String("value2.0"),
	)

	require.Equal(t, 2, len(set))

	newSet := set.WithTags(
		Key("key1").String("value1.1"),
		Key("key3").String("value3.0"),
	)

	require.Equal(t, 3, len(newSet))
	assert.Equal(t, "value1.1", newSet[0].Value.Interface)
	assert.Equal(t, "value2.0", newSet[1].Value.Interface)
	assert.Equal(t, "value3.0", newSet[2].Value.Interface)

	require.Len(t, set, 2)
	assert.Equal(t, "value1.0", set[0].Value.Interface)
	assert.Equal(t, "value2.0", set[1].Value.Interface)
}

func TestWithTags_Chained(t *testing.T) {
	t.Run("chains when new tag has Chained(true)", func(t *testing.T) {
		set := EmptySet.WithTags(Key("key1").String("a"))
		newSet := set.WithTags(Key("key1").String("b").Chained(true))

		require.Len(t, newSet, 1)
		assert.Equal(t, "a.b", newSet[0].Value.Interface)
		require.Len(t, set, 1)
		assert.Equal(t, "a", set[0].Value.Interface)
	})

	t.Run("replaces when new tag has no Chained flag", func(t *testing.T) {
		set := EmptySet.WithTags(Key("key1").String("a"))
		newSet := set.WithTags(Key("key1").String("b"))

		require.Len(t, newSet, 1)
		assert.Equal(t, "b", newSet[0].Value.Interface)
	})

	t.Run("chains repeatedly for the same key", func(t *testing.T) {
		set := EmptySet.WithTags(Key("key1").String("a"))
		set = set.WithTags(Key("key1").String("b").Chained(true))
		set = set.WithTags(Key("key1").String("c").Chained(true))

		require.Len(t, set, 1)
		assert.Equal(t, "a.b.c", set[0].Value.Interface)
	})

	t.Run("new tag can interrupt a chain by not setting Chained", func(t *testing.T) {
		set := EmptySet.WithTags(Key("key1").String("a"))
		set = set.WithTags(Key("key1").String("b").Chained(true))
		set = set.WithTags(Key("key1").String("fresh"))

		require.Len(t, set, 1)
		assert.Equal(t, "fresh", set[0].Value.Interface)
	})

	t.Run("original set is not mutated", func(t *testing.T) {
		set := EmptySet.WithTags(Key("key1").String("a"))
		_ = set.WithTags(Key("key1").String("b").Chained(true))

		require.Len(t, set, 1)
		assert.Equal(t, "a", set[0].Value.Interface)
	})

	t.Run("other keys unchanged when chaining one key", func(t *testing.T) {
		set := EmptySet.WithTags(
			Key("key1").String("x"),
			Key("key2").String("y"),
		)
		newSet := set.WithTags(Key("key1").String("z").Chained(true))

		require.Len(t, newSet, 2)
		assert.Equal(t, "x.z", newSet[0].Value.Interface)
		assert.Equal(t, "y", newSet[1].Value.Interface)
	})
}

func TestMapTag(t *testing.T) {
	t.Run("Map value holds a tag set", func(t *testing.T) {
		m := Key("meta").Map(
			Key("a").String("1"),
			Key("b").String("2"),
		)
		require.Equal(t, MAP, m.Value.Type)
		inner := m.Value.Interface.(Set)
		require.Len(t, inner, 2)
		assert.Equal(t, "1", tagStringInSet(inner, "a"))
		assert.Equal(t, "2", tagStringInSet(inner, "b"))
	})

	t.Run("chained maps merge disjoint inner keys", func(t *testing.T) {
		old := EmptySet.WithTags(Key("meta").Map(Key("a").String("1")))
		merged := old.WithTags(Key("meta").Map(Key("b").String("2")).Chained(true))

		require.Len(t, merged, 1)
		inner := asMapSet(t, merged[0])
		require.Len(t, inner, 2)
		assert.Equal(t, "1", tagStringInSet(inner, "a"))
		assert.Equal(t, "2", tagStringInSet(inner, "b"))
	})

	t.Run("chained maps merge same inner key for strings", func(t *testing.T) {
		old := EmptySet.WithTags(Key("meta").Map(Key("x").String("a")))
		// Inner tag for x must be Chained to concatenate; outer meta must be Chained to merge maps.
		merged := old.WithTags(Key("meta").Map(Key("x").String("b").Chained(true)).Chained(true))

		inner := asMapSet(t, merged[0])
		require.Len(t, inner, 1)
		assert.Equal(t, "a.b", tagStringInSet(inner, "x"))
	})

	t.Run("nested map merge chains inner map entries", func(t *testing.T) {
		old := EmptySet.WithTags(Key("outer").Map(
			Key("inner").Map(Key("k").String("v1")),
		))
		merged := old.WithTags(Key("outer").Map(
			Key("inner").Map(Key("k").String("v2").Chained(true)).Chained(true),
		).Chained(true))

		require.Len(t, merged, 1)
		outer := asMapSet(t, merged[0])
		require.Len(t, outer, 1)
		innerTag := outer[0]
		require.Equal(t, MAP, innerTag.Value.Type)
		inner := asMapSet(t, innerTag)
		require.Len(t, inner, 1)
		assert.Equal(t, "v1.v2", tagStringInSet(inner, "k"))
	})

	t.Run("nested merge three levels", func(t *testing.T) {
		old := EmptySet.WithTags(Key("l1").Map(
			Key("l2").Map(
				Key("l3").Map(Key("leaf").String("a")),
			),
		))
		merged := old.WithTags(Key("l1").Map(
			Key("l2").Map(
				Key("l3").Map(Key("leaf").String("b").Chained(true)).Chained(true),
			).Chained(true),
		).Chained(true))

		l1 := asMapSet(t, merged[0])
		l2 := asMapSet(t, l1[0])
		l3 := asMapSet(t, l2[0])
		assert.Equal(t, "a.b", tagStringInSet(l3, "leaf"))
	})

	t.Run("chained map merges inner map with extra keys at each layer", func(t *testing.T) {
		old := EmptySet.WithTags(Key("meta").Map(
			Key("keep").String("old"),
			Key("shared").String("s1"),
		))
		merged := old.WithTags(Key("meta").Map(
			Key("shared").String("s2").Chained(true),
			Key("extra").String("new"),
		).Chained(true))

		inner := asMapSet(t, merged[0])
		require.Len(t, inner, 3)
		assert.Equal(t, "old", tagStringInSet(inner, "keep"))
		assert.Equal(t, "s1.s2", tagStringInSet(inner, "shared"))
		assert.Equal(t, "new", tagStringInSet(inner, "extra"))
	})

	t.Run("replace map when new tag is not chained", func(t *testing.T) {
		old := EmptySet.WithTags(Key("meta").Map(Key("a").String("1")))
		merged := old.WithTags(Key("meta").Map(Key("b").String("2")))

		inner := asMapSet(t, merged[0])
		require.Len(t, inner, 1)
		assert.Equal(t, "2", tagStringInSet(inner, "b"))
	})

	t.Run("chained string on old map new string panics", func(t *testing.T) {
		old := EmptySet.WithTags(Key("meta").Map(Key("a").String("1")))
		merged := old.WithTags(Key("meta").String("b").Chained(true))
		assert.Equal(t, "b", merged[0].Value.Interface)
	})
}

func asMapSet(t *testing.T, tag *Tag) Set {
	t.Helper()
	require.Equal(t, MAP, tag.Value.Type)
	return tag.Value.Interface.(Set)
}

func tagStringInSet(s Set, key Key) string {
	for _, tg := range s {
		if tg.Key == key && tg.Value.Type == STRING {
			return tg.Value.Interface.(string)
		}
	}
	return ""
}
