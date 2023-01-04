package logs

import (
	"fmt"
	"github.com/kuoss/lethe/storage/driver/factory"
	"os"
	"path/filepath"
	"testing"

	"github.com/kuoss/lethe/testutil"
)

func Test_GetDiskUsedKB(t *testing.T) {

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}

	kb, err := rotator.GetDiskUsedKB(rotator.driver.RootDirectory())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("kb=", kb)
}

func Test_DeleteByAge(t *testing.T) {
	testutil.SetTestLogFiles()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp2", "log")})
	rotator := Rotator{driver: d}

	rotator.DeleteByAge(false)

	want := `[{"Filepath":"/tmp/log/node/node01/2009-11-10_20.log","LogType":"node","Target":"node01","Name":"2009-11-10_20","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_21.log","LogType":"node","Target":"node01","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_22.log","LogType":"node","Target":"node01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_23.log","LogType":"node","Target":"node01","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_20.log","LogType":"node","Target":"node02","Name":"2009-11-10_20","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_21.log","LogType":"node","Target":"node02","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_22.log","LogType":"node","Target":"node02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_23.log","LogType":"node","Target":"node02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_20.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_20","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_21.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_22.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_23.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_20.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_20","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_21.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_22.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_23.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_23","Extention":"log","KB":0}]`
	got := rotator.ListFiles()

	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got), "16")
}

func Test_DeleteBySize(t *testing.T) {
	testutil.SetTestLogFiles()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp2", "log")})
	rotator := Rotator{driver: d}

	rotator.DeleteBySize(false)
	want := `[{"Filepath":"/tmp/log/node/node01/2009-11-10_22.log","LogType":"node","Target":"node01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_23.log","LogType":"node","Target":"node01","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_22.log","LogType":"node","Target":"node02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_23.log","LogType":"node","Target":"node02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_21.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_22.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_23.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_22.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_23.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_23","Extention":"log","KB":0}]`
	got := rotator.ListFiles()

	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got), "9")
}

func Test_DeleteBySize_DryRun(t *testing.T) {
	testutil.SetTestLogFiles()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}

	rotator.DeleteBySize(true)
}
