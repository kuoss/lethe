package main

import (
	"testing"

	"github.com/kuoss/lethe/testutil"
)

func Test_list_dirs(t *testing.T) {

	var got string
	var want string

	testutil.SetTestLogFiles()
	got = execute("list", "dirs")
	want = `"DIR                        SIZE(Mi)   FILES   FIRST               LAST              \n/tmp/log/node/node01            6.0       6   2009-11-10_18.log   2009-11-10_23.log   \n/tmp/log/node/node02            6.0       6   2009-11-10_18.log   2009-11-10_23.log   \n/tmp/log/pod/namespace02        6.0       6   2009-11-10_18.log   2009-11-10_23.log   \n/tmp/log/pod/namspace01         6.0       6   2009-11-10_18.log   2009-11-10_23.log   \nTOTAL                          24.0      24   -                   -"`
	testutil.CheckEqualJSON(t, got, want)

	testutil.SetTestLogFiles()
	got = execute("list", "files")
	want = `"FILEPATH                                     SIZE(Mi) \n/tmp/log/node/node01/2009-11-10_18.log            1.0   \n/tmp/log/node/node01/2009-11-10_19.log            1.0   \n/tmp/log/node/node01/2009-11-10_20.log            1.0   \n/tmp/log/node/node01/2009-11-10_21.log            1.0   \n/tmp/log/node/node01/2009-11-10_22.log            1.0   \n/tmp/log/node/node01/2009-11-10_23.log            1.0   \n/tmp/log/node/node02/2009-11-10_18.log            1.0   \n/tmp/log/node/node02/2009-11-10_19.log            1.0   \n/tmp/log/node/node02/2009-11-10_20.log            1.0   \n/tmp/log/node/node02/2009-11-10_21.log            1.0   \n/tmp/log/node/node02/2009-11-10_22.log            1.0   \n/tmp/log/node/node02/2009-11-10_23.log            1.0   \n/tmp/log/pod/namespace02/2009-11-10_18.log        1.0   \n/tmp/log/pod/namespace02/2009-11-10_19.log        1.0   \n/tmp/log/pod/namespace02/2009-11-10_20.log        1.0   \n/tmp/log/pod/namespace02/2009-11-10_21.log        1.0   \n/tmp/log/pod/namespace02/2009-11-10_22.log        1.0   \n/tmp/log/pod/namespace02/2009-11-10_23.log        1.0   \n/tmp/log/pod/namspace01/2009-11-10_18.log         1.0   \n/tmp/log/pod/namspace01/2009-11-10_19.log         1.0   \n/tmp/log/pod/namspace01/2009-11-10_20.log         1.0   \n/tmp/log/pod/namspace01/2009-11-10_21.log         1.0   \n/tmp/log/pod/namspace01/2009-11-10_22.log         1.0   \n/tmp/log/pod/namspace01/2009-11-10_23.log         1.0   \nTOTAL                                            24.0"`
	testutil.CheckEqualJSON(t, got, want)
}
