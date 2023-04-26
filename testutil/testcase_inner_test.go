package testutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTC_inner(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{TC(), "TESTCASE:testcase_inner_test.go:15"},
		{TC(), "TESTCASE:testcase_inner_test.go:16"},
		{TC(), "TESTCASE:testcase_inner_test.go:17"},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			assert.Equal(t, tc.want, tc.input)
		})
	}
}
