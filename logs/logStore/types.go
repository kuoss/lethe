package logStore

import "fmt"

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
	CompactRaw() string
}

type NodeLog struct {
	Name    string `json:"-"`
	Time    string `json:"time,omitempty"`
	Node    string `json:"node,omitempty"`
	Process string `json:"process,omitempty"`
	Log     string `json:"log,omitempty"`
}

func (log NodeLog) CompactRaw() string {
	return fmt.Sprintf("%s[%s|%s] %s", log.Time, log.Node, log.Process, log.Log)
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
	Time      string `json:"time"`
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
	Container string `json:"container"`
	Log       string `json:"log"`
}

func (log PodLog) CompactRaw() string {
	return fmt.Sprintf("%s[%s|%s|%s] %s", log.Time, log.Namespace, log.Pod, log.Container, log.Log)
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
