package main

import (
	"testing"

	"github.com/kuoss/lethe/testutil"
)

func Test_task_deleteByAge(t *testing.T) {
	testutil.SetTestLogFiles()
	got := execute("task", "delete-by-age")
	want := `"DeleteByAge(2009-11-10_18 \u003c 2009-11-10_20): /tmp/log/node/node02/2009-11-10_18.log\nDeleteByAge(2009-11-10_18 \u003c 2009-11-10_20): /tmp/log/pod/namspace01/2009-11-10_18.log\nDeleteByAge(2009-11-10_18 \u003c 2009-11-10_20): /tmp/log/pod/namespace02/2009-11-10_18.log\nDeleteByAge(2009-11-10_18 \u003c 2009-11-10_20): /tmp/log/node/node01/2009-11-10_18.log\nDeleteByAge(2009-11-10_19 \u003c 2009-11-10_20): /tmp/log/node/node01/2009-11-10_19.log\nDeleteByAge(2009-11-10_19 \u003c 2009-11-10_20): /tmp/log/pod/namspace01/2009-11-10_19.log\nDeleteByAge(2009-11-10_19 \u003c 2009-11-10_20): /tmp/log/node/node02/2009-11-10_19.log\nDeleteByAge(2009-11-10_19 \u003c 2009-11-10_20): /tmp/log/pod/namespace02/2009-11-10_19.log\nDeleteByAge(2009-11-10_20): done"`
	testutil.CheckEqualJSON(t, got, want)
}

func Test_task_deleteByAge_dryRun(t *testing.T) {
	testutil.SetTestLogFiles()
	got := execute("task", "delete-by-age", "--dry-run")
	want := `"DeleteByAge(2009-11-10_18 \u003c 2009-11-10_20): /tmp/log/node/node02/2009-11-10_18.log (dry run)\nDeleteByAge(2009-11-10_18 \u003c 2009-11-10_20): /tmp/log/pod/namspace01/2009-11-10_18.log (dry run)\nDeleteByAge(2009-11-10_18 \u003c 2009-11-10_20): /tmp/log/pod/namespace02/2009-11-10_18.log (dry run)\nDeleteByAge(2009-11-10_18 \u003c 2009-11-10_20): /tmp/log/node/node01/2009-11-10_18.log (dry run)\nDeleteByAge(2009-11-10_19 \u003c 2009-11-10_20): /tmp/log/node/node01/2009-11-10_19.log (dry run)\nDeleteByAge(2009-11-10_19 \u003c 2009-11-10_20): /tmp/log/pod/namspace01/2009-11-10_19.log (dry run)\nDeleteByAge(2009-11-10_19 \u003c 2009-11-10_20): /tmp/log/node/node02/2009-11-10_19.log (dry run)\nDeleteByAge(2009-11-10_19 \u003c 2009-11-10_20): /tmp/log/pod/namespace02/2009-11-10_19.log (dry run)\nDeleteByAge(2009-11-10_20): done"`
	testutil.CheckEqualJSON(t, got, want)
}

func Test_task_deleteBySize(t *testing.T) {
	testutil.SetTestLogFiles()
	got := execute("task", "delete-by-size")
	want := `"DeleteBySize(24.0m \u003e 10.0m): /tmp/log/node/node02/2009-11-10_18.log\nDeleteBySize(23.0m \u003e 10.0m): /tmp/log/pod/namspace01/2009-11-10_18.log\nDeleteBySize(22.0m \u003e 10.0m): /tmp/log/pod/namespace02/2009-11-10_18.log\nDeleteBySize(21.0m \u003e 10.0m): /tmp/log/node/node01/2009-11-10_18.log\nDeleteBySize(20.0m \u003e 10.0m): /tmp/log/node/node01/2009-11-10_19.log\nDeleteBySize(19.0m \u003e 10.0m): /tmp/log/pod/namspace01/2009-11-10_19.log\nDeleteBySize(18.0m \u003e 10.0m): /tmp/log/node/node02/2009-11-10_19.log\nDeleteBySize(17.0m \u003e 10.0m): /tmp/log/pod/namespace02/2009-11-10_19.log\nDeleteBySize(16.0m \u003e 10.0m): /tmp/log/node/node01/2009-11-10_20.log\nDeleteBySize(15.0m \u003e 10.0m): /tmp/log/node/node02/2009-11-10_20.log\nDeleteBySize(14.0m \u003e 10.0m): /tmp/log/pod/namspace01/2009-11-10_20.log\nDeleteBySize(13.0m \u003e 10.0m): /tmp/log/pod/namespace02/2009-11-10_20.log\nDeleteBySize(12.0m \u003e 10.0m): /tmp/log/node/node02/2009-11-10_21.log\nDeleteBySize(11.0m \u003e 10.0m): /tmp/log/pod/namspace01/2009-11-10_21.log\nDeleteBySize(10.0m \u003e 10.0m): /tmp/log/node/node01/2009-11-10_21.log\nDeleteBySize(9.0m \u003c 10.0m): done"`
	testutil.CheckEqualJSON(t, got, want)
}

func Test_task_deleteBySize_dryRun(t *testing.T) {
	testutil.SetTestLogFiles()
	got := execute("task", "delete-by-size", "--dry-run")
	want := `"DeleteBySize(24.0m \u003e 10.0m): /tmp/log/node/node02/2009-11-10_18.log (dry run)\nDeleteBySize(23.0m \u003e 10.0m): /tmp/log/pod/namspace01/2009-11-10_18.log (dry run)\nDeleteBySize(22.0m \u003e 10.0m): /tmp/log/pod/namespace02/2009-11-10_18.log (dry run)\nDeleteBySize(21.0m \u003e 10.0m): /tmp/log/node/node01/2009-11-10_18.log (dry run)\nDeleteBySize(20.0m \u003e 10.0m): /tmp/log/node/node01/2009-11-10_19.log (dry run)\nDeleteBySize(19.0m \u003e 10.0m): /tmp/log/pod/namspace01/2009-11-10_19.log (dry run)\nDeleteBySize(18.0m \u003e 10.0m): /tmp/log/node/node02/2009-11-10_19.log (dry run)\nDeleteBySize(17.0m \u003e 10.0m): /tmp/log/pod/namespace02/2009-11-10_19.log (dry run)\nDeleteBySize(16.0m \u003e 10.0m): /tmp/log/node/node01/2009-11-10_20.log (dry run)\nDeleteBySize(15.0m \u003e 10.0m): /tmp/log/node/node02/2009-11-10_20.log (dry run)\nDeleteBySize(14.0m \u003e 10.0m): /tmp/log/pod/namspace01/2009-11-10_20.log (dry run)\nDeleteBySize(13.0m \u003e 10.0m): /tmp/log/pod/namespace02/2009-11-10_20.log (dry run)\nDeleteBySize(12.0m \u003e 10.0m): /tmp/log/node/node02/2009-11-10_21.log (dry run)\nDeleteBySize(11.0m \u003e 10.0m): /tmp/log/pod/namspace01/2009-11-10_21.log (dry run)\nDeleteBySize(10.0m \u003e 10.0m): /tmp/log/node/node01/2009-11-10_21.log (dry run)\nDeleteBySize(9.0m \u003c 10.0m): done"`
	testutil.CheckEqualJSON(t, got, want)
}
