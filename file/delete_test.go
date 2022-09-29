package file

import (
	"fmt"
	"testing"

	"github.com/kuoss/lethe/testutil"
)

func Test_GetDiskUsedKB(t *testing.T) {
	kb, err := GetDiskUsedKB()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("kb=", kb)
}

func Test_DeleteByAge(t *testing.T) {
	testutil.SetTestLogFiles()
	DeleteByAge(false)
	want := `[{"Filepath":"/tmp/log/node/node01/2009-11-10_20.log","Typ":"node","Target":"node01","Name":"2009-11-10_20","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_21.log","Typ":"node","Target":"node01","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_22.log","Typ":"node","Target":"node01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_23.log","Typ":"node","Target":"node01","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_20.log","Typ":"node","Target":"node02","Name":"2009-11-10_20","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_21.log","Typ":"node","Target":"node02","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_22.log","Typ":"node","Target":"node02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_23.log","Typ":"node","Target":"node02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_20.log","Typ":"pod","Target":"namespace02","Name":"2009-11-10_20","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_21.log","Typ":"pod","Target":"namespace02","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_22.log","Typ":"pod","Target":"namespace02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_23.log","Typ":"pod","Target":"namespace02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_20.log","Typ":"pod","Target":"namspace01","Name":"2009-11-10_20","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_21.log","Typ":"pod","Target":"namspace01","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_22.log","Typ":"pod","Target":"namspace01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_23.log","Typ":"pod","Target":"namspace01","Name":"2009-11-10_23","Extention":"log","KB":0}]`
	got := ListFiles()
	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got), "16")
}

func Test_DeleteBySize(t *testing.T) {
	testutil.SetTestLogFiles()
	DeleteBySize(false)
	want := `[{"Filepath":"/tmp/log/node/node01/2009-11-10_22.log","Typ":"node","Target":"node01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node01/2009-11-10_23.log","Typ":"node","Target":"node01","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_22.log","Typ":"node","Target":"node02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/node/node02/2009-11-10_23.log","Typ":"node","Target":"node02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_21.log","Typ":"pod","Target":"namespace02","Name":"2009-11-10_21","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_22.log","Typ":"pod","Target":"namespace02","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namespace02/2009-11-10_23.log","Typ":"pod","Target":"namespace02","Name":"2009-11-10_23","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_22.log","Typ":"pod","Target":"namspace01","Name":"2009-11-10_22","Extention":"log","KB":0},{"Filepath":"/tmp/log/pod/namspace01/2009-11-10_23.log","Typ":"pod","Target":"namspace01","Name":"2009-11-10_23","Extention":"log","KB":0}]`
	got := ListFiles()
	testutil.CheckEqualJSON(t, got, want)
	testutil.CheckEqualJSON(t, len(got), "9")
}

func Test_DeleteBySize_DryRun(t *testing.T) {
	testutil.SetTestLogFiles()
	DeleteBySize(true)
}
