package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kuoss/lethe/storage"
	"github.com/kuoss/lethe/storage/driver/factory"
)

func TestDepth(t *testing.T) {
	userHomeDirectory, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHomeDirectory, "tmp", "log")})
	logPath := storage.LogPath{RootDirectory: d.RootDirectory()}
	logPath.SetFullPath(filepath.Join("pod", "namespace01", "2022-11-10_23.log"))
	rtn := logPath.Depth()
	fmt.Println(rtn)
}
