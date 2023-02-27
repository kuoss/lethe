package logs

import (
	"testing"

	"github.com/kuoss/lethe/config"

	"github.com/kuoss/lethe/storage/driver/factory"

	"github.com/kuoss/lethe/testutil"
)

func TestListFiles(t *testing.T) {
	testutil.SetTestLogFiles()

	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": config.GetLogRoot()})
	rotator := Rotator{driver: d}
	got := rotator.ListFiles()
	want := `[{"FullPath":"tmp/log/node/node01/2009-11-10_21.log","SubPath":"2009-11-10_21.log","LogType":"node","Target":"node01","Name":"2009-11-10_21.log","Extention":".log","Size":1057},{"FullPath":"tmp/log/node/node01/2009-11-10_22.log","SubPath":"2009-11-10_22.log","LogType":"node","Target":"node01","Name":"2009-11-10_22.log","Extention":".log","Size":177},{"FullPath":"tmp/log/node/node02/2009-11-01_00.log","SubPath":"2009-11-01_00.log","LogType":"node","Target":"node02","Name":"2009-11-01_00.log","Extention":".log","Size":0},{"FullPath":"tmp/log/node/node02/2009-11-10_21.log","SubPath":"2009-11-10_21.log","LogType":"node","Target":"node02","Name":"2009-11-10_21.log","Extention":".log","Size":1116},{"FullPath":"tmp/log/pod/namespace01/2000-01-01_00.log","SubPath":"2000-01-01_00.log","LogType":"pod","Target":"namespace01","Name":"2000-01-01_00.log","Extention":".log","Size":1031},{"FullPath":"tmp/log/pod/namespace01/2009-11-10_21.log","SubPath":"2009-11-10_21.log","LogType":"pod","Target":"namespace01","Name":"2009-11-10_21.log","Extention":".log","Size":279},{"FullPath":"tmp/log/pod/namespace01/2009-11-10_22.log","SubPath":"2009-11-10_22.log","LogType":"pod","Target":"namespace01","Name":"2009-11-10_22.log","Extention":".log","Size":1031},{"FullPath":"tmp/log/pod/namespace01/2029-11-10_23.log","SubPath":"2029-11-10_23.log","LogType":"pod","Target":"namespace01","Name":"2029-11-10_23.log","Extention":".log","Size":279},{"FullPath":"tmp/log/pod/namespace02/0000-00-00_00.log","SubPath":"0000-00-00_00.log","LogType":"pod","Target":"namespace02","Name":"0000-00-00_00.log","Extention":".log","Size":12},{"FullPath":"tmp/log/pod/namespace02/2009-11-10_22.log","SubPath":"2009-11-10_22.log","LogType":"pod","Target":"namespace02","Name":"2009-11-10_22.log","Extention":".log","Size":1125}]`
	testutil.CheckEqualJSON(t, got, want)
}

func TestListDirs(t *testing.T) {
	testutil.SetTestLogFiles()

	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": config.GetLogRoot()})
	rotator := Rotator{driver: d}
	got := rotator.ListDirs()
	want := `[{"FullPath":"tmp/log/node/node01","SubPath":"node/node01","LogType":"node","Target":"node01","FileCount":0,"FirstFile":"","LastFile":"","Size":0,"LastForward":""},{"FullPath":"tmp/log/node/node02","SubPath":"node/node02","LogType":"node","Target":"node02","FileCount":0,"FirstFile":"","LastFile":"","Size":0,"LastForward":""},{"FullPath":"tmp/log/pod/namespace01","SubPath":"pod/namespace01","LogType":"pod","Target":"namespace01","FileCount":0,"FirstFile":"","LastFile":"","Size":0,"LastForward":""},{"FullPath":"tmp/log/pod/namespace02","SubPath":"pod/namespace02","LogType":"pod","Target":"namespace02","FileCount":0,"FirstFile":"","LastFile":"","Size":0,"LastForward":""}]`
	testutil.CheckEqualJSON(t, got, want)
}

func TestListDirWithSize(t *testing.T) {
	testutil.SetTestLogFiles()

	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": config.GetLogRoot()})
	rotator := Rotator{driver: d}
	got := rotator.ListDirsWithSize()

	want := `[{"FullPath":"tmp/log/node/node01","SubPath":"node/node01","LogType":"node","Target":"node01","FileCount":2,"FirstFile":"2009-11-10_21.log","LastFile":"2009-11-10_22.log","Size":1234,"LastForward":""},{"FullPath":"tmp/log/node/node02","SubPath":"node/node02","LogType":"node","Target":"node02","FileCount":2,"FirstFile":"2009-11-01_00.log","LastFile":"2009-11-10_21.log","Size":1116,"LastForward":""},{"FullPath":"tmp/log/pod/namespace01","SubPath":"pod/namespace01","LogType":"pod","Target":"namespace01","FileCount":4,"FirstFile":"2000-01-01_00.log","LastFile":"2029-11-10_23.log","Size":2620,"LastForward":""},{"FullPath":"tmp/log/pod/namespace02","SubPath":"pod/namespace02","LogType":"pod","Target":"namespace02","FileCount":2,"FirstFile":"0000-00-00_00.log","LastFile":"2009-11-10_22.log","Size":1137,"LastForward":""}]`
	testutil.CheckEqualJSON(t, got, want)
}

func TestListTargets(t *testing.T) {
	testutil.SetTestLogFiles()

	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": config.GetLogRoot()})
	rotator := Rotator{driver: d}
	got := rotator.ListTargets()

	want := `[{"FullPath":"tmp/log/node/node01","SubPath":"node/node01","LogType":"node","Target":"node01","FileCount":2,"FirstFile":"2009-11-10_21.log","LastFile":"2009-11-10_22.log","Size":1234,"LastForward":"2009-11-10T23:00:00."},{"FullPath":"tmp/log/node/node02","SubPath":"node/node02","LogType":"node","Target":"node02","FileCount":2,"FirstFile":"2009-11-01_00.log","LastFile":"2009-11-10_21.log","Size":1116,"LastForward":"2009-11-10T21:58:00."},{"FullPath":"tmp/log/pod/namespace01","SubPath":"pod/namespace01","LogType":"pod","Target":"namespace01","FileCount":4,"FirstFile":"2000-01-01_00.log","LastFile":"2029-11-10_23.log","Size":2620,"LastForward":"2009-11-10T23:00:00."},{"FullPath":"tmp/log/pod/namespace02","SubPath":"pod/namespace02","LogType":"pod","Target":"namespace02","FileCount":2,"FirstFile":"0000-00-00_00.log","LastFile":"2009-11-10_22.log","Size":1137,"LastForward":"2009-11-10T22:58:00."}]`
	testutil.CheckEqualJSON(t, got, want)
}
