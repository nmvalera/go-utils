package mock

import (
	"testing"

	store "github.com/nmvalera/go-utils/store"
	"github.com/stretchr/testify/assert"
)

func TestImplementsStore(t *testing.T) {
	assert.Implements(t, (*store.Store)(nil), new(MockStore))
}
