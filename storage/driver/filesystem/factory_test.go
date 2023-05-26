package filesystem

import (
	"testing"

	"github.com/kuoss/lethe/storage/driver/factory"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	want := New(Params{RootDirectory: "asdf"})
	got, err := factory.Get("filesystem", map[string]interface{}{"RootDirectory": "asdf"})
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
