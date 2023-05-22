package filesystem

import (
	"bufio"
	"os"
	"testing"

	storagedriver "github.com/kuoss/lethe/storage/driver"
	"github.com/kuoss/lethe/util/testutil"
	"github.com/stretchr/testify/assert"
)

func init() {
	dir := "tmp/storage_driver_filesystem_filewriter_test"
	testutil.ChdirRoot()
	err := os.RemoveAll(dir)
	if err != nil {
		panic(err)
	}
	err = os.Mkdir(dir, 0755)
	if err != nil {
		panic(err)
	}
}

func TestFileWriter(t *testing.T) {
	f, err := os.Create("tmp/storage_driver_filesystem_filewriter_test/greet.txt")
	assert.NoError(t, err)
	bw := bufio.NewWriter(f)

	var fw storagedriver.FileWriter = &fileWriter{
		file:      f,
		size:      0,
		bw:        bw,
		closed:    false,
		committed: false,
		cancelled: false,
	}
	s := "hello"
	n, err := fw.Write([]byte(s))
	assert.NoError(t, err)
	assert.Equal(t, len(s), n)
	fw.Close()

	content, err := os.ReadFile("tmp/storage_driver_filesystem_filewriter_test/greet.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(content))
}
