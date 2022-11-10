package routine

import (
	"testing"
	"time"

	"github.com/kuoss/lethe/file"
	"github.com/kuoss/lethe/testutil"
)

func Test_RoutineDelete(t *testing.T) {
	testutil.Init()
	testutil.SetTestLogFiles()
	Start(time.Duration(1))
	time.Sleep(3 * time.Second)
	got := file.ListFiles()
	want := `[{"Filepath":"/tmp/log/node/node01/2009-11-10_22.log","Typ":"node","Target":"node01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_23.log","Typ":"node","Target":"node01","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_22.log","Typ":"node","Target":"node02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_23.log","Typ":"node","Target":"node02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_21.log","Typ":"pod","Target":"namespace02","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_22.log","Typ":"pod","Target":"namespace02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_23.log","Typ":"pod","Target":"namespace02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_22.log","Typ":"pod","Target":"namspace01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_23.log","Typ":"pod","Target":"namspace01","Name":"2009-11-10_23","Extention":"log","KB":0}]`
	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got), "9")
}
