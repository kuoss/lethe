package logStore

const (
	AUDIT_TYPE = "audit"
	EVENT_TYPE = "event"
	POD_TYPE   = "pod"
	NODE_TYPE  = "node"
)

type LogType interface {
	GetName() string
	process() int
}

type NodeLog struct {
	Name string
}

func (log NodeLog) process() int {
	return 0
}

func (log NodeLog) GetName() string {
	return log.Name
}

type PodLog struct {
	Name string
}

func (log PodLog) process() int {

	return 0
}

func (log PodLog) GetName() string {
	return log.Name
}
