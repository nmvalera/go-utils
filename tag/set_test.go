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
