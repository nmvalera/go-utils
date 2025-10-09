package common //no:revive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrt(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		res := Ptr("test")
		assert.Equal(t, "test", *res)
	})
	t.Run("int", func(t *testing.T) {
		res := Ptr(1)
		assert.Equal(t, 1, *res)
	})
	t.Run("bool", func(t *testing.T) {
		res := Ptr(true)
		assert.Equal(t, true, *res)
	})
}

func TestVal(t *testing.T) {
	t.Run("string#nil", func(t *testing.T) {
		var v *string
		res := Val(v)
		assert.Equal(t, "", res)
	})
	t.Run("string#non-il", func(t *testing.T) {
		v := "test"
		res := Val(&v)
		assert.Equal(t, "test", res)
	})
	t.Run("int#nil", func(t *testing.T) {
		var v *int
		res := Val(v)
		assert.Equal(t, 0, res)
	})
	t.Run("int#non-nil", func(t *testing.T) {
		v := 1
		res := Val(&v)
		assert.Equal(t, 1, res)
	})
	t.Run("bool#nil", func(t *testing.T) {
		var v *bool
		res := Val(v)
		assert.Equal(t, false, res)
	})
	t.Run("bool#non-nil", func(t *testing.T) {
		v := true
		res := Val(&v)
		assert.Equal(t, true, res)
	})
}

func TestPtrSlice(t *testing.T) {
	t.Run("NonEmpty", func(t *testing.T) {
		res := PtrSlice("test", "test2")
		assert.Len(t, *res, 2)
		assert.Equal(t, "test", *(*res)[0])
		assert.Equal(t, "test2", *(*res)[1])
	})

	t.Run("Empty", func(t *testing.T) {
		res := PtrSlice[string]()
		assert.Equal(t, []*string{}, *res)
	})
}

func TestValSlice(t *testing.T) {
	t.Run("NonEmpty", func(t *testing.T) {
		res := ValSlice(Ptr("test"), Ptr("test2"))
		assert.Equal(t, []string{"test", "test2"}, res)
	})
	t.Run("Empty", func(t *testing.T) {
		res := ValSlice[string]()
		assert.Equal(t, []string{}, res)
	})
}
