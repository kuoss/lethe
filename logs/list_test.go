package logs

import (
	"fmt"
	"github.com/kuoss/lethe/storage/driver/factory"
	"os"
	"path/filepath"
	"testing"

	"github.com/kuoss/lethe/testutil"
)

func Test_List(t *testing.T) {
	testutil.SetTestLogs()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}

	got := rotator.ListFiles()
	want := `[{"Filepath":"/tmp/log/node/node01/2009-11-10_22.log","LogType":"node","Target":"node01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_23.log","LogType":"node","Target":"node01","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_22.log","LogType":"node","Target":"node02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace01/2009-11-10_22.log","LogType":"pod","Target":"namespace01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace01/2009-11-10_23.log","LogType":"pod","Target":"namespace01","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_22.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_22","Extention":"log","KB":0}]`
	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got), "6")
}

func TestListDirs(t *testing.T) {
	testutil.SetTestLogs()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}

	got := rotator.ListDirs()
	for _, logDir := range got {
		fmt.Printf("logDir : %+v\n", logDir)
	}
}

func TestListDirWithSize(t *testing.T) {
	testutil.SetTestLogs()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}

	got := rotator.ListDirsWithSize()
	for _, logDir := range got {
		fmt.Printf("logDir : %+v\n", logDir)
	}
}

func TestListTargets(t *testing.T) {
	testutil.SetTestLogs()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}

	got := rotator.ListTargets()
	for _, logDir := range got {
		fmt.Printf("logDir : %+v\n", logDir)
	}
}
