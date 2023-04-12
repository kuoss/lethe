package logStore

const (
	AUDIT_TYPE = "audit"
	EVENT_TYPE = "event"
	POD_TYPE   = "pod"
	NODE_TYPE  = "node"
)

type LogLine interface {
	GetName() string
	process() int
	getTime() string
}

type NodeLog struct {
	Name    string `json:"-"`
	Time    string `json:"time,omitempty"`
	Node    string `json:"node,omitempty"`
	Process string `json:"process,omitempty"`
	Log     string `json:"log,omitempty"`
}

func (log NodeLog) process() int {
	return 0
}

func (log NodeLog) GetName() string {
	return log.Name
}

func (log NodeLog) getTime() string {
	return log.Time
}

type PodLog struct {
	Name      string `json:"-"`
	Time      string `json:"time,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Pod       string `json:"pod,omitempty"`
	Container string `json:"container,omitempty"`
	Log       string `json:"log,omitempty"`
}

func (log PodLog) process() int {

	return 0
}

func (log PodLog) GetName() string {
	return log.Name
}

func (log PodLog) getTime() string {
	return log.Time
}
