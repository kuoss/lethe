package file

import (
	"testing"

	"github.com/kuoss/lethe/testutil"
)

func Test_List(t *testing.T) {
	testutil.SetTestLogs()

	var got interface{}
	var want string

	got = ListFiles()
	want = `[{"Filepath":"/tmp/log/node/node01/2009-11-10_22.log","Typ":"node","Target":"node01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_23.log","Typ":"node","Target":"node01","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_22.log","Typ":"node","Target":"node02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace01/2009-11-10_22.log","Typ":"pod","Target":"namespace01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace01/2009-11-10_23.log","Typ":"pod","Target":"namespace01","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_22.log","Typ":"pod","Target":"namespace02","Name":"2009-11-10_22","Extention":"log","KB":0}]`
	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got.([]LogFile)), "6")

	got = ListFilesWithSize()
	want = `[{"Filepath":"/tmp/log/node/node01/2009-11-10_22.log","Typ":"node","Target":"node01","Name":"2009-11-10_22","Extention":"log","KB":4},{"Filepath":"/tmp/log/node/node01/2009-11-10_23.log","Typ":"node","Target":"node01","Name":"2009-11-10_23","Extention":"log","KB":4},{"Filepath":"/tmp/log/node/node02/2009-11-10_22.log","Typ":"node","Target":"node02","Name":"2009-11-10_22","Extention":"log","KB":4},{"Filepath":"/tmp/log/pod/namespace01/2009-11-10_22.log","Typ":"pod","Target":"namespace01","Name":"2009-11-10_22","Extention":"log","KB":4},{"Filepath":"/tmp/log/pod/namespace01/2009-11-10_23.log","Typ":"pod","Target":"namespace01","Name":"2009-11-10_23","Extention":"log","KB":4},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_22.log","Typ":"pod","Target":"namespace02","Name":"2009-11-10_22","Extention":"log","KB":4}]`
	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got.([]LogFile)), "6")

	got = ListDirs()
	want = `["node/node01","node/node02","pod/namespace01","pod/namespace02"]`
	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got.([]string)), "4")

	got = ListTargets()
	want = `[{"Dirpath":"/tmp/log/node/node01","Typ":"node","Target":"node01","CountFiles":2,"FirstFile":"2009-11-10_22.log","LastFile":"2009-11-10_23.log","KB":8,"LastLog":"2009-11-10T23:02:00.000000Z"},{"Dirpath":"/tmp/log/node/node02","Typ":"node","Target":"node02","CountFiles":1,"FirstFile":"2009-11-10_22.log","LastFile":"2009-11-10_22.log","KB":4,"LastLog":"2009-11-10T22:58:00.000000Z"},{"Dirpath":"/tmp/log/pod/namespace01","Typ":"pod","Target":"namespace01","CountFiles":2,"FirstFile":"2009-11-10_22.log","LastFile":"2009-11-10_23.log","KB":8,"LastLog":"2009-11-10T23:02:00.000000Z"},{"Dirpath":"/tmp/log/pod/namespace02","Typ":"pod","Target":"namespace02","CountFiles":1,"FirstFile":"2009-11-10_22.log","LastFile":"2009-11-10_22.log","KB":4,"LastLog":"2009-11-10T22:58:00.000000Z"}]`
	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got.([]LogDir)), "4")
}
