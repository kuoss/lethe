package logmodel

import "fmt"

type LogType string

const (
	LogTypeNode LogType = "node"
	LogTypePod  LogType = "pod"
)

type NodeLog struct {
	Time    string `json:"time"`
	Node    string `json:"node"`
	Process string `json:"process"`
	Log     string `json:"log"`
}

func (l NodeLog) Type() LogType  { return LogTypeNode }
func (l NodeLog) String() string { return fmt.Sprintf("%#v", l) }

type PodLog struct {
	Time      string `json:"time"`
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
	Container string `json:"container"`
	Log       string `json:"log"`
}

func (l PodLog) Type() LogType  { return LogTypePod }
func (l PodLog) String() string { return fmt.Sprintf("%#v", l) }
