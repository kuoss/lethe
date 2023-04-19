package rotator

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/kuoss/lethe/testutil"
	"github.com/stretchr/testify/assert"
)

func TestListDirs(t *testing.T) {
	testutil.SetTestLogFiles()

	got := NewRotator().ListDirs()

	actual := fmt.Sprintf("%v", got)
	expected := "[{tmp/log/node/node01 node/node01 node node01 0   0 } {tmp/log/node/node02 node/node02 node node02 0   0 } {tmp/log/pod/namespace01 pod/namespace01 pod namespace01 0   0 } {tmp/log/pod/namespace02 pod/namespace02 pod namespace02 0   0 }]"
	if runtime.GOOS == "windows" {
		expected = "[{tmp\\log\\node\\node01 node\\node01 node node01 0   0 } {tmp\\log\\node\\node02 node\\node02 node node02 0   0 } {tmp\\log\\pod\\namespace01 pod\\namespace01 pod namespace01 0   0 } {tmp\\log\\pod\\namespace02 pod\\namespace02 pod namespace02 0   0 }]"
	}
	assert.Equal(t, expected, actual)
}

func TestListDirWithSize(t *testing.T) {
	testutil.SetTestLogFiles()

	got := NewRotator().ListDirsWithSize()

	actual := fmt.Sprintf("%v", got)
	expected := "[{tmp/log/node/node01 node/node01 node node01 2 2009-11-10_21.log 2009-11-10_22.log 1234 } {tmp/log/node/node02 node/node02 node node02 2 2009-11-01_00.log 2009-11-10_21.log 1116 } {tmp/log/pod/namespace01 pod/namespace01 pod namespace01 4 2000-01-01_00.log 2029-11-10_23.log 2620 } {tmp/log/pod/namespace02 pod/namespace02 pod namespace02 2 0000-00-00_00.log 2009-11-10_22.log 1137 }]"
	if runtime.GOOS == "windows" {
		expected = "[{tmp\\log\\node\\node01 node\\node01 node node01 2 2009-11-10_21.log 2009-11-10_22.log 1248 } {tmp\\log\\node\\node02 node\\node02 node node02 2 2009-11-01_00.log 2009-11-10_21.log 1128 } {tmp\\log\\pod\\namespace01 pod\\namespace01 pod namespace01 4 2000-01-01_00.log 2029-11-10_23.log 2646 } {tmp\\log\\pod\\namespace02 pod\\namespace02 pod namespace02 2 0000-00-00_00.log 2009-11-10_22.log 1151 }]"
	}
	assert.Equal(t, expected, actual)
}

func TestListFiles(t *testing.T) {
	testutil.SetTestLogFiles()

	got := NewRotator().ListFiles()

	actual := fmt.Sprintf("%v", got)
	expected := "[{tmp/log/node/node01/2009-11-10_21.log 2009-11-10_21.log node node01 2009-11-10_21.log .log 1057} {tmp/log/node/node01/2009-11-10_22.log 2009-11-10_22.log node node01 2009-11-10_22.log .log 177} {tmp/log/node/node02/2009-11-01_00.log 2009-11-01_00.log node node02 2009-11-01_00.log .log 0} {tmp/log/node/node02/2009-11-10_21.log 2009-11-10_21.log node node02 2009-11-10_21.log .log 1116} {tmp/log/pod/namespace01/2000-01-01_00.log 2000-01-01_00.log pod namespace01 2000-01-01_00.log .log 1031} {tmp/log/pod/namespace01/2009-11-10_21.log 2009-11-10_21.log pod namespace01 2009-11-10_21.log .log 279} {tmp/log/pod/namespace01/2009-11-10_22.log 2009-11-10_22.log pod namespace01 2009-11-10_22.log .log 1031} {tmp/log/pod/namespace01/2029-11-10_23.log 2029-11-10_23.log pod namespace01 2029-11-10_23.log .log 279} {tmp/log/pod/namespace02/0000-00-00_00.log 0000-00-00_00.log pod namespace02 0000-00-00_00.log .log 12} {tmp/log/pod/namespace02/2009-11-10_22.log 2009-11-10_22.log pod namespace02 2009-11-10_22.log .log 1125}]"
	if runtime.GOOS == "windows" {
		expected = "[{tmp\\log\\node\\node01\\2009-11-10_21.log 2009-11-10_21.log node node01 2009-11-10_21.log .log 1068} {tmp\\log\\node\\node01\\2009-11-10_22.log 2009-11-10_22.log node node01 2009-11-10_22.log .log 180} {tmp\\log\\node\\node02\\2009-11-01_00.log 2009-11-01_00.log node node02 2009-11-01_00.log .log 0} {tmp\\log\\node\\node02\\2009-11-10_21.log 2009-11-10_21.log node node02 2009-11-10_21.log .log 1128} {tmp\\log\\pod\\namespace01\\2000-01-01_00.log 2000-01-01_00.log pod namespace01 2000-01-01_00.log .log 1041} {tmp\\log\\pod\\namespace01\\2009-11-10_21.log 2009-11-10_21.log pod namespace01 2009-11-10_21.log .log 282} {tmp\\log\\pod\\namespace01\\2009-11-10_22.log 2009-11-10_22.log pod namespace01 2009-11-10_22.log .log 1041} {tmp\\log\\pod\\namespace01\\2029-11-10_23.log 2029-11-10_23.log pod namespace01 2029-11-10_23.log .log 282} {tmp\\log\\pod\\namespace02\\0000-00-00_00.log 0000-00-00_00.log pod namespace02 0000-00-00_00.log .log 14} {tmp\\log\\pod\\namespace02\\2009-11-10_22.log 2009-11-10_22.log pod namespace02 2009-11-10_22.log .log 1137}]"
	}
	assert.Equal(t, expected, actual)
}

func TestListTargets(t *testing.T) {
	testutil.SetTestLogFiles()

	got := NewRotator().ListTargets()

	actual := fmt.Sprintf("%v", got)
	expected := "[{tmp/log/node/node01 node/node01 node node01 2 2009-11-10_21.log 2009-11-10_22.log 1234 2009-11-10T23:00:00.} {tmp/log/node/node02 node/node02 node node02 2 2009-11-01_00.log 2009-11-10_21.log 1116 2009-11-10T21:58:00.} {tmp/log/pod/namespace01 pod/namespace01 pod namespace01 4 2000-01-01_00.log 2029-11-10_23.log 2620 2009-11-10T23:00:00.} {tmp/log/pod/namespace02 pod/namespace02 pod namespace02 2 0000-00-00_00.log 2009-11-10_22.log 1137 2009-11-10T22:58:00.}]"
	if runtime.GOOS == "windows" {
		expected = "[{tmp\\log\\node\\node01 node\\node01 node node01 2 2009-11-10_21.log 2009-11-10_22.log 1248 2009-11-10T23:00:00.} {tmp\\log\\node\\node02 node\\node02 node node02 2 2009-11-01_00.log 2009-11-10_21.log 1128 2009-11-10T21:58:00.} {tmp\\log\\pod\\namespace01 pod\\namespace01 pod namespace01 4 2000-01-01_00.log 2029-11-10_23.log 2646 2009-11-10T23:00:00.} {tmp\\log\\pod\\namespace02 pod\\namespace02 pod namespace02 2 0000-00-00_00.log 2009-11-10_22.log 1151 2009-11-10T22:58:00.}]"
	}
	assert.Equal(t, expected, actual)
}
