package logs

import (
	"github.com/kuoss/lethe/storage/driver/factory"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kuoss/lethe/testutil"
)

func Test_RoutineDelete(t *testing.T) {
	testutil.Init()
	testutil.SetTestLogFiles()

	userHome, _ := os.UserHomeDir()
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": filepath.Join(userHome, "tmp", "log")})
	rotator := Rotator{driver: d}

	rotator.Start(time.Duration(1))
	time.Sleep(3 * time.Second)
	got := rotator.ListFiles()
	want := `[{"Filepath":"/tmp/log/node/node01/2009-11-10_22.log","LogType":"node","Target":"node01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_23.log","LogType":"node","Target":"node01","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_22.log","LogType":"node","Target":"node02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_23.log","LogType":"node","Target":"node02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_21.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_22.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_23.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_22.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_23.log","LogType":"pod","Target":"namspace01","Name":"2009-11-10_23","Extention":"log","KB":0}]`
	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got), "9")
}
