package testutil_test

import (
	"fmt"
	"testing"

	"github.com/kuoss/lethe/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTC_outer(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{testutil.TC(), "TESTCASE:testcase_outer_test.go:16"},
		{testutil.TC(), "TESTCASE:testcase_outer_test.go:17"},
		{testutil.TC(), "TESTCASE:testcase_outer_test.go:18"},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			assert.Equal(t, tc.want, tc.input)
		})
	}
}
