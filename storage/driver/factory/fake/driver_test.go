package fake

import (
	"fmt"
	"io"
	"testing"

	storagedriver "github.com/kuoss/lethe/storage/driver"
	"github.com/stretchr/testify/assert"
)

var (
	driver1 = New()
)

func TestNew(t *testing.T) {
	assert.NotNil(t, driver1)
	assert.Equal(t, "*fake.driver", fmt.Sprintf("%T", driver1))
}

func TestName(t *testing.T) {
	got := driver1.Name()
	assert.Equal(t, "fake", got)
}

func TestGetContent(t *testing.T) {
	got, err := driver1.GetContent("")
	assert.NoError(t, err)
	assert.Equal(t, []byte{}, got)
}

func TestPutContent(t *testing.T) {
	err := driver1.PutContent("", []byte{})
	assert.NoError(t, err)
}

func TestReader(t *testing.T) {
	got, err := driver1.Reader("")
	assert.NoError(t, err)
	assert.Equal(t, &io.PipeReader{}, got)
}

func TestStat(t *testing.T) {
	got, err := driver1.Stat("")
	assert.NoError(t, err)
	assert.Equal(t, nil, got)
}

func TestList(t *testing.T) {
	got, err := driver1.List("")
	assert.NoError(t, err)
	assert.Equal(t, []string{}, got)
}

func TestMove(t *testing.T) {
	err := driver1.Move("", "")
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	err := driver1.Delete("")
	assert.NoError(t, err)
}

func TestWalk(t *testing.T) {
	got, err := driver1.Walk("")
	assert.NoError(t, err)
	assert.Equal(t, ([]storagedriver.FileInfo)(nil), got)
}

func TestWalkDir(t *testing.T) {
	got, err := driver1.WalkDir("")
	assert.NoError(t, err)
	assert.Equal(t, []string{}, got)
}

func TestWriter(t *testing.T) {
	got, err := driver1.Writer("")
	assert.NoError(t, err)
	assert.Equal(t, nil, got)
}

func TestRootDirectory(t *testing.T) {
	got := driver1.RootDirectory()
	assert.Equal(t, "", got)
}
