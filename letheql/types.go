package letheql

type ParsedQuery struct {
	Type    string
	Labels  []Label
	Keyword string
}

type Label struct {
	Key   string
	Value string
}

type ValueType string

const (
	ValueTypeScalar ValueType = "scalar"
	ValueTypeLogs   ValueType = "logs"
)

type QueryData struct {
	ResultType ValueType `json:"resultType"`
	Logs       []string  `json:"logs,omitempty"`
	Scalar     float64   `json:"scalar,omitempty"`
}
