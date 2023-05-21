package main

import (
	"runtime"
	"testing"

	"github.com/kuoss/lethe/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_list_dirs(t *testing.T) {
	testutil.SetTestLogFiles()

	actual := execute("list", "dirs")

	expected := "DIR                       SIZE(Mi)   FILES   FIRST FILE          LAST FILE         \ntmp/log/node/node01            0.0       2   2009-11-10_21.log   2009-11-10_22.log   \ntmp/log/node/node02            0.0       2   2009-11-01_00.log   2009-11-10_21.log   \ntmp/log/pod/namespace01        0.0       4   2000-01-01_00.log   2029-11-10_23.log   \ntmp/log/pod/namespace02        0.0       2   0000-00-00_00.log   2009-11-10_22.log   \nTOTAL                          0.0      10   -                   -"
	if runtime.GOOS == "windows" {
		expected = "DIR                       SIZE(Mi)   FILES   FIRST FILE          LAST FILE         \ntmp\\log\\node\\node01            0.0       2   2009-11-10_21.log   2009-11-10_22.log   \ntmp\\log\\node\\node02            0.0       2   2009-11-01_00.log   2009-11-10_21.log   \ntmp\\log\\pod\\namespace01        0.0       4   2000-01-01_00.log   2029-11-10_23.log   \ntmp\\log\\pod\\namespace02        0.0       2   0000-00-00_00.log   2009-11-10_22.log   \nTOTAL                          0.0      10   -                   -"
	}
	assert.Equal(t, expected, actual)
}

func Test_list_files(t *testing.T) {
	testutil.SetTestLogFiles()

	actual := execute("list", "files")
	expected := "FILEPATH                                    SIZE(Mi) \ntmp/log/node/node01/2009-11-10_21.log            0.0   \ntmp/log/node/node01/2009-11-10_22.log            0.0   \ntmp/log/node/node02/2009-11-01_00.log            0.0   \ntmp/log/node/node02/2009-11-10_21.log            0.0   \ntmp/log/pod/namespace01/2000-01-01_00.log        0.0   \ntmp/log/pod/namespace01/2009-11-10_21.log        0.0   \ntmp/log/pod/namespace01/2009-11-10_22.log        0.0   \ntmp/log/pod/namespace01/2029-11-10_23.log        0.0   \ntmp/log/pod/namespace02/0000-00-00_00.log        0.0   \ntmp/log/pod/namespace02/2009-11-10_22.log        0.0   \nTOTAL                                            0.0"
	if runtime.GOOS == "windows" {
		expected = "FILEPATH                                    SIZE(Mi) \ntmp\\log\\node\\node01\\2009-11-10_21.log            0.0   \ntmp\\log\\node\\node01\\2009-11-10_22.log            0.0   \ntmp\\log\\node\\node02\\2009-11-01_00.log            0.0   \ntmp\\log\\node\\node02\\2009-11-10_21.log            0.0   \ntmp\\log\\pod\\namespace01\\2000-01-01_00.log        0.0   \ntmp\\log\\pod\\namespace01\\2009-11-10_21.log        0.0   \ntmp\\log\\pod\\namespace01\\2009-11-10_22.log        0.0   \ntmp\\log\\pod\\namespace01\\2029-11-10_23.log        0.0   \ntmp\\log\\pod\\namespace02\\0000-00-00_00.log        0.0   \ntmp\\log\\pod\\namespace02\\2009-11-10_22.log        0.0   \nTOTAL                                            0.0"
	}
	assert.Equal(t, expected, actual)
}
