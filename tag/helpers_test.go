package tag

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sortTagsByKey(tags []*Tag) {
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Key < tags[j].Key
	})
}

func TestMapTags(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		tags := MapTags(map[string]any{"k": "v"})
		require.Len(t, tags, 1)
		assert.Equal(t, Key("k"), tags[0].Key)
		assert.Equal(t, STRING, tags[0].Value.Type)
		assert.Equal(t, "v", tags[0].Value.Interface)
	})

	t.Run("int64", func(t *testing.T) {
		tags := MapTags(map[string]any{"n": int64(42)})
		require.Len(t, tags, 1)
		assert.Equal(t, INT64, tags[0].Value.Type)
		assert.Equal(t, int64(42), tags[0].Value.Interface)
	})

	t.Run("int", func(t *testing.T) {
		tags := MapTags(map[string]any{"n": 7})
		require.Len(t, tags, 1)
		assert.Equal(t, INT64, tags[0].Value.Type)
		assert.Equal(t, int64(7), tags[0].Value.Interface)
	})

	t.Run("float64", func(t *testing.T) {
		tags := MapTags(map[string]any{"f": 2.5})
		require.Len(t, tags, 1)
		assert.Equal(t, FLOAT64, tags[0].Value.Type)
		assert.Equal(t, 2.5, tags[0].Value.Interface)
	})

	t.Run("bool", func(t *testing.T) {
		tags := MapTags(map[string]any{"b": true})
		require.Len(t, tags, 1)
		assert.Equal(t, BOOL, tags[0].Value.Type)
		assert.Equal(t, true, tags[0].Value.Interface)
	})

	t.Run("nested map[string]any", func(t *testing.T) {
		tags := MapTags(map[string]any{
			"outer": map[string]any{
				"inner": "x",
				"n":     int(3),
			},
		})
		require.Len(t, tags, 1)
		assert.Equal(t, Key("outer"), tags[0].Key)
		require.Equal(t, MAP, tags[0].Value.Type)
		inner := tags[0].Value.Interface.(Set)
		sortTagsByKey(inner)
		require.Len(t, inner, 2)
		assert.Equal(t, STRING, inner[0].Value.Type)
		assert.Equal(t, "x", inner[0].Value.Interface)
		assert.Equal(t, Key("inner"), inner[0].Key)
		assert.Equal(t, Key("n"), inner[1].Key)
		assert.Equal(t, INT64, inner[1].Value.Type)
		assert.Equal(t, int64(3), inner[1].Value.Interface)
	})

	t.Run("deeply nested maps", func(t *testing.T) {
		tags := MapTags(map[string]any{
			"a": map[string]any{
				"b": map[string]any{
					"c": "leaf",
				},
			},
		})
		require.Len(t, tags, 1)
		l1 := tags[0].Value.Interface.(Set)
		require.Len(t, l1, 1)
		l2 := l1[0].Value.Interface.(Set)
		require.Len(t, l2, 1)
		leaf := l2[0]
		assert.Equal(t, Key("c"), leaf.Key)
		assert.Equal(t, STRING, leaf.Value.Type)
		assert.Equal(t, "leaf", leaf.Value.Interface)
	})

	t.Run("default falls back to object", func(t *testing.T) {
		type custom struct{ X string }
		v := custom{X: "z"}
		tags := MapTags(map[string]any{"obj": v})
		require.Len(t, tags, 1)
		assert.Equal(t, OBJECT, tags[0].Value.Type)
		assert.Equal(t, v, tags[0].Value.Interface)
	})

	t.Run("multiple top-level keys stable assertion", func(t *testing.T) {
		tags := MapTags(map[string]any{
			"z": "last",
			"a": "first",
			"m": 1,
		})
		require.Len(t, tags, 3)
		sortTagsByKey(tags)
		assert.Equal(t, Key("a"), tags[0].Key)
		assert.Equal(t, "first", tags[0].Value.Interface)
		assert.Equal(t, Key("m"), tags[1].Key)
		assert.Equal(t, int64(1), tags[1].Value.Interface)
		assert.Equal(t, Key("z"), tags[2].Key)
		assert.Equal(t, "last", tags[2].Value.Interface)
	})

	t.Run("empty map", func(t *testing.T) {
		assert.Empty(t, MapTags(map[string]any{}))
	})
}
