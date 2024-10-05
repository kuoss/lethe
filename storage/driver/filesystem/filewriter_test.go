package filesystem

import (
	"bufio"
	"os"
	"testing"

	"github.com/kuoss/common/tester"
	storagedriver "github.com/kuoss/lethe/storage/driver"
	"github.com/stretchr/testify/require"
)

func TestFileWriter(t *testing.T) {
	_, cleanup := tester.SetupDir(t, map[string]string{
		"@/testdata/log": "data/log",
	})
	defer cleanup()

	f, err := os.Create("data/log/greet.txt")
	require.NoError(t, err)
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
	require.NoError(t, err)
	require.Equal(t, len(s), n)
	fw.Close()

	content, err := os.ReadFile("data/log/greet.txt")
	require.NoError(t, err)
	require.Equal(t, "hello", string(content))
}
