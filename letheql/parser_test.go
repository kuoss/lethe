package letheql

import (
	"testing"

	"github.com/kuoss/lethe/testutil"
)

func Test_ParseQuerySuccess(t *testing.T) {
	var query string
	var want string
	var got ParsedQuery

	want = `{"Type":"pod","Labels":null,"Keyword":""}`
	got, _ = ParseQuery(`pod`)
	testutil.CheckEqualJSON(t, got, want)

	want = `{"Type":"pod","Labels":null,"Keyword":""}`
	got, _ = ParseQuery(`pod{}`)
	testutil.CheckEqualJSON(t, got, want)

	query = `pod{} hello`
	want = `{"Type":"pod","Labels":null,"Keyword":"hello"}`
	got, _ = ParseQuery(query)
	testutil.CheckEqualJSON(t, got, want, "query=", query)

	query = `pod{namespace="namespace01",pod="nginx-*"} hello`
	want = `{"Type":"pod","Labels":[{"Key":"namespace","Value":"namespace01"},{"Key":"pod","Value":"nginx-*"}],"Keyword":"hello"}`
	got, _ = ParseQuery(query)
	testutil.CheckEqualJSON(t, got, want, "query=", query)

	query = `node{node="node01",process="kubelet"} hello`
	want = `{"Type":"node","Labels":[{"Key":"node","Value":"node01"},{"Key":"process","Value":"kubelet"}],"Keyword":"hello"}`
	got, _ = ParseQuery(query)
	testutil.CheckEqualJSON(t, got, want, "query=", query)

	query = `pod{namespace=namespace01,pod=nginx-*}`
	want = `{"Type":"pod","Labels":[{"Key":"namespace","Value":"namespace01"},{"Key":"pod","Value":"nginx-*"}],"Keyword":""}`
	got, _ = ParseQuery(query)
	testutil.CheckEqualJSON(t, got, want, "query=", query)
}

func Test_ParseQueryFail(t *testing.T) {
	var want string
	var got error

	want = `"invalid query type"`
	_, got = ParseQuery("hello")
	testutil.CheckEqualJSON(t, got, want)
}
