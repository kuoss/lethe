package parser

import (
	"errors"
	"github.com/kuoss/lethe/letheql"
	"testing"

	"github.com/kuoss/lethe/testutil"
	"github.com/stretchr/testify/assert"
)

func init() {
	testutil.Init()
}

func Test_ParseQuerySuccess(t *testing.T) {

	tests := map[string]struct {
		query string
		want  letheql.ParsedQuery
	}{
		`query pod`:         {query: `pod`, want: letheql.ParsedQuery{Type: "pod", Labels: nil, Keyword: ""}},
		`query pod{}`:       {query: `pod{}`, want: letheql.ParsedQuery{Type: "pod", Labels: nil, Keyword: ""}},
		`query pod{} hello`: {query: `pod{} hello`, want: letheql.ParsedQuery{Type: "pod", Labels: nil, Keyword: "hello"}},
		`query pod{namespace="namespace01",pod="nginx-*"} hello`: {query: `pod{namespace="namespace01",pod="nginx-*"} hello`, want: letheql.ParsedQuery{Type: "pod", Labels: []letheql.Label{{Key: "namespace", Value: "namespace01"}, {Key: "pod", Value: "nginx-*"}}, Keyword: "hello"}},
		`query node{node="node01",process="kubelet"} hello`:      {query: `node{node="node01",process="kubelet"} hello`, want: letheql.ParsedQuery{Type: "node", Labels: []letheql.Label{{Key: "node", Value: "node01"}, {Key: "process", Value: "kubelet"}}, Keyword: "hello"}},
		`query pod{namespace="namespace01",pod="nginx-*"}`:       {query: `pod{namespace="namespace01",pod="nginx-*"}`, want: letheql.ParsedQuery{Type: "pod", Labels: []letheql.Label{{Key: "namespace", Value: "namespace01"}, {Key: "pod", Value: "nginx-*"}}, Keyword: ""}},
	}
	for name, tt := range tests {
		t.Run(name, func(subt *testing.T) {
			got, err := ParseQuery(tt.query)
			if err != nil {
				subt.Fatalf("query: %s err: %s", name, err.Error())
			}
			assert.Equal(subt, tt.want, got)
		})
	}
}

func Test_ParseQueryFail(t *testing.T) {

	want := errors.New("invalid query type")
	_, err := ParseQuery("hello")
	if assert.Error(t, err) {
		assert.Equal(t, want, err)
	}
}
