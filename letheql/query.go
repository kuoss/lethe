package letheql

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/file"
	"github.com/VictoriaMetrics/metricsql"
)

type LeafType string

const (
	LeafTypeAuditLogRequest LeafType = "AuditLogRequest"
	LeafTypeEventLogRequest LeafType = "EventLogRequest"
	LeafTypeNodeLogRequest  LeafType = "NodeLogRequest"
	LeafTypePodLogRequest   LeafType = "PodLogRequest"
)

const (
	LeafTypeScalarResult LeafType = "ScalarResult"
	LeafTypeLogsResult   LeafType = "LogsResult"
)

type Leaf struct {
	LeafType     LeafType
	LogRequest   LogRequest
	ScalarResult float64
	LogsResult   []string
	Function     string
	TimeRange    TimeRange
	Keywords     []string
}

type LogRequest struct {
	// AuditSearchParams file.AuditSearchParams
	EventSearchParams file.EventSearchParams
	NodeSearchParams  file.NodeSearchParams
	PodSearchParams   file.PodSearchParams
	DurationSeconds   int
	Function          string
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func ProcQuery(query string, timeRange TimeRange) (QueryData, error) {
	fmt.Println("==> ProcQuery", query, "timeRange=", timeRange)
	parts := strings.Split(query, "|")
	if len(parts) < 1 {
		return QueryData{}, errors.New("empty query")
	}
	parsableQuery := parts[0]
	keywords := parts[1:]
	for i, v := range keywords {
		keywords[i] = strings.TrimSpace(v)
	}
	expr, err := metricsql.Parse(parsableQuery)
	if err != nil {
		log.Println("ProcQuery: Parse: err=", err)
		return QueryData{}, err
	}
	leaf, err := procExpr(expr, Leaf{TimeRange: timeRange, Keywords: keywords})
	if err != nil {
		log.Println("ProcQuery: procExpr: err=", err)
		return QueryData{}, err
	}
	leaf, err = resolveLeaf(leaf)
	if err != nil {
		log.Println("ProcQuery: resolveLeaf: err=", err)
		return QueryData{}, err
	}

	queryData, err := getQueryDataFromLeaf(leaf)
	// fmt.Printf("ProcQuery: leaf=%#v\n", leaf)
	// fmt.Printf("ProcQuery: queryData=%#v\n", queryData)
	if err != nil {
		log.Println("ProcQuery: getQueryDataFromLeaf: err=", err)
		return QueryData{}, err
	}
	return queryData, nil
}

func getQueryDataFromLeaf(leaf Leaf) (QueryData, error) {
	var queryData QueryData
	// fmt.Printf("getQueryDataFromLeaf: leaf=%#v\n", leaf)
	switch leaf.LeafType {
	case LeafTypeLogsResult:
		queryData.ResultType = ValueTypeLogs
		queryData.Logs = leaf.LogsResult
	case LeafTypeScalarResult:
		queryData.ResultType = ValueTypeScalar
		queryData.Scalar = leaf.ScalarResult
	default:
		return queryData, errors.New("log request not resolved")
	}
	return queryData, nil
}

func resolveLeaf(leaf Leaf) (Leaf, error) {
	if leaf.LeafType != LeafTypeAuditLogRequest &&
		leaf.LeafType != LeafTypeEventLogRequest &&
		leaf.LeafType != LeafTypeNodeLogRequest &&
		leaf.LeafType != LeafTypePodLogRequest {
		return leaf, nil
	}
	req := leaf.LogRequest
	// DurationSeconds, TimeRange{} => DurationSeconds, EndTime
	now := config.GetNow()
	if leaf.TimeRange.End.IsZero() {
		leaf.TimeRange.End = now
	}
	if leaf.TimeRange.Start.IsZero() {
		leaf.TimeRange.Start = leaf.TimeRange.End.Add(time.Duration(-40*24) * time.Hour)
	}
	if leaf.TimeRange.End == leaf.TimeRange.Start {
		return Leaf{}, errors.New("end time and start time are the same")
	}
	if leaf.TimeRange.End.Before(leaf.TimeRange.Start) {
		return Leaf{}, errors.New("end time is earlier than start time")
	}
	durationSecondsFromTimeRange := int(leaf.TimeRange.End.Sub(leaf.TimeRange.Start) / 1000 / 1000 / 1000)
	if durationSecondsFromTimeRange > 40*86400 {
		durationSecondsFromTimeRange = 40 * 86400
	}
	durationSeconds := req.DurationSeconds
	// fmt.Println("leaf.TimeRange.Start=", leaf.TimeRange.Start, "leaf.TimeRange.End=", leaf.TimeRange.End)
	// fmt.Println("durationSeconds=", durationSeconds, "durationSecondsFromTimeRange=", durationSecondsFromTimeRange)
	if durationSeconds == 0 || (durationSecondsFromTimeRange != 0 && durationSecondsFromTimeRange < durationSeconds) {
		durationSeconds = durationSecondsFromTimeRange
	}
	logSearch := file.LogSearch{
		DurationSeconds: durationSeconds,
		EndTime:         leaf.TimeRange.End,
		Keywords:        leaf.Keywords,
	}
	switch leaf.LeafType {
	case LeafTypeAuditLogRequest:
		logSearch.LogType = "audit"
		logSearch.TargetPattern = "audit"
		// logSearch.AuditSearchParams = req.AuditSearchParams
	case LeafTypeEventLogRequest:
		logSearch.LogType = "event"
		logSearch.TargetPattern = "event"
		logSearch.EventSearchParams = req.EventSearchParams
	case LeafTypeNodeLogRequest:
		logSearch.LogType = "node"
		logSearch.TargetPattern = req.NodeSearchParams.NodePattern
		logSearch.NodeSearchParams = req.NodeSearchParams
	case LeafTypePodLogRequest:
		logSearch.LogType = "pod"
		logSearch.TargetPattern = req.PodSearchParams.NamespacePattern
		logSearch.PodSearchParams = req.PodSearchParams
	}
	switch req.Function {
	case "":
		result, err := file.GetLogs(logSearch)
		if err != nil {
			return Leaf{}, err
		}
		leaf.LogsResult = result.Logs
		leaf.LeafType = LeafTypeLogsResult
	case "count_over_time":
		logSearch.IsCounting = true
		result, err := file.GetLogs(logSearch)
		if err != nil {
			return Leaf{}, err
		}
		leaf.ScalarResult = float64(result.Count)
		leaf.LeafType = LeafTypeScalarResult
	default:
		return leaf, fmt.Errorf("not supported function: %s", req.Function)
	}
	return leaf, nil
}

func procExpr(expr metricsql.Expr, leaf Leaf) (Leaf, error) {
	// fmt.Printf("procExpr: %#T %#v\n", expr, leaf.TimeRange)
	switch v := expr.(type) {
	case *metricsql.BinaryOpExpr:
		return procBinaryOpExpr(v, leaf)
	case *metricsql.FuncExpr:
		return procFuncExpr(v, leaf)
	case *metricsql.MetricExpr:
		return procMetricExpr(v, leaf)
	case *metricsql.NumberExpr:
		return procNumberExpr(v, leaf)
	case *metricsql.RollupExpr:
		return procRollupExpr(v, leaf)
	}
	return leaf, nil
}

func procFuncExpr(expr *metricsql.FuncExpr, leaf Leaf) (Leaf, error) {
	leaf.Function = expr.Name
	newLeaf := Leaf{}
	for _, arg := range expr.Args {
		var err error
		newLeaf, err = procExpr(arg, leaf)
		if err != nil {
			return Leaf{}, err
		}
		// fmt.Printf("procFuncExpr: newLeaf=%#v\n", newLeaf)
	}
	return newLeaf, nil
}

func procBinaryOpExpr(expr *metricsql.BinaryOpExpr, leaf Leaf) (Leaf, error) {
	// TODO: should be vector not scalar
	var leftLeaf, rightLeaf Leaf
	var err error
	leftLeaf, err = procExpr(expr.Left, leaf)
	if err != nil {
		return Leaf{}, err
	}
	rightLeaf, err = procExpr(expr.Right, leaf)
	if err != nil {
		return Leaf{}, err
	}
	leftLeaf, err = resolveLeaf(leftLeaf)
	if err != nil {
		return Leaf{}, err
	}
	rightLeaf, err = resolveLeaf(rightLeaf)
	if err != nil {
		return Leaf{}, err
	}
	if leftLeaf.LeafType != LeafTypeScalarResult || rightLeaf.LeafType != LeafTypeScalarResult {
		return Leaf{}, errors.New("not allowed leafType for operator")
	}
	leaf.LeafType = LeafTypeScalarResult
	leaf.ScalarResult = 0
	switch expr.Op {
	case ">":
		if leftLeaf.ScalarResult > rightLeaf.ScalarResult {
			leaf.ScalarResult = 1
		}
	case "<":
		if leftLeaf.ScalarResult < rightLeaf.ScalarResult {
			leaf.ScalarResult = 1
		}
	case "==":
		if leftLeaf.ScalarResult == rightLeaf.ScalarResult {
			leaf.ScalarResult = 1
		}
	case "!=":
		if leftLeaf.ScalarResult != rightLeaf.ScalarResult {
			leaf.ScalarResult = 1
		}
	default:
		return leaf, fmt.Errorf("not supported operator: %s", expr.Op)
	}
	return leaf, nil
}

func procMetricExpr(expr *metricsql.MetricExpr, leaf Leaf) (Leaf, error) {
	if len(expr.LabelFilters) < 1 {
		return Leaf{}, errors.New("must have one or more labels")
	}
	// fmt.Printf("expr=[%#v] LabelFilters=[%#v]\n", expr.LabelFilters, expr)
	if expr.LabelFilters[0].Label != "__name__" {
		return Leaf{}, errors.New("a log name must be specified")
	}
	switch expr.LabelFilters[0].Value {
	case "audit":
		return procAuditExpr(expr, leaf)
	case "event":
		return procEventExpr(expr, leaf)
	case "node":
		return procNodeExpr(expr, leaf)
	case "pod":
		return procPodExpr(expr, leaf)
	}
	return Leaf{}, errors.New("unknown log name")
}

func procAuditExpr(expr *metricsql.MetricExpr, leaf Leaf) (Leaf, error) {
	leaf.LeafType = LeafTypeAuditLogRequest
	leaf.LogRequest = LogRequest{}
	return leaf, nil
}

func procEventExpr(expr *metricsql.MetricExpr, leaf Leaf) (Leaf, error) {
	var namespace, typ, reason, object, count string
	for _, l := range expr.LabelFilters {
		switch l.Label {
		case "namespace":
			namespace = l.Value
		case "type":
			typ = l.Value
		case "reason":
			reason = l.Value
		case "object":
			object = l.Value
		case "count":
			count = l.Value
		}
	}
	leaf.LeafType = LeafTypeEventLogRequest
	leaf.LogRequest = LogRequest{EventSearchParams: file.EventSearchParams{Namespace: namespace, Type: typ, Reason: reason, Object: object, Count: count}, Function: leaf.Function}
	return leaf, nil
}

func procNodeExpr(expr *metricsql.MetricExpr, leaf Leaf) (Leaf, error) {
	var node, process string
	for _, l := range expr.LabelFilters {
		switch l.Label {
		case "node":
			node = l.Value
		case "process":
			process = l.Value
		}
	}
	leaf.LeafType = LeafTypeNodeLogRequest
	leaf.LogRequest = LogRequest{NodeSearchParams: file.NodeSearchParams{NodePattern: node, Process: process}, Function: leaf.Function}
	return leaf, nil
}

func procPodExpr(expr *metricsql.MetricExpr, leaf Leaf) (Leaf, error) {
	var namespace, pod, container string
	for _, l := range expr.LabelFilters {
		switch l.Label {
		case "namespace":
			// explicitly empty (label exists) => error
			if l.Value == "" {
				return Leaf{}, errors.New("namespace value cannot be empty")
			}
			namespace = l.Value
		case "pod":
			pod = l.Value
		case "container":
			container = l.Value
		case "__name__":
		default:
			return Leaf{}, errors.New("unknown label " + l.Label)
		}
	}
	// implicit empty (label not exists) => all
	if namespace == "" {
		namespace = "*"
	}
	leaf.LeafType = LeafTypePodLogRequest
	leaf.LogRequest = LogRequest{PodSearchParams: file.PodSearchParams{NamespacePattern: namespace, Pod: pod, Container: container}, Function: leaf.Function}
	return leaf, nil

}
func procNumberExpr(expr *metricsql.NumberExpr, leaf Leaf) (Leaf, error) {
	leaf.LeafType = LeafTypeScalarResult
	leaf.ScalarResult = expr.N
	return leaf, nil
}

func procRollupExpr(expr *metricsql.RollupExpr, leaf Leaf) (Leaf, error) {
	if reflect.ValueOf(expr.Window).Type().String() != "*metricsql.DurationExpr" {
		return Leaf{}, errors.New("not duration expr")
	}
	leaf, err := procExpr(expr.Expr, leaf)
	if err != nil {
		return Leaf{}, err
	}
	if leaf.LeafType == LeafTypePodLogRequest || leaf.LeafType == LeafTypeNodeLogRequest || leaf.LeafType == LeafTypeEventLogRequest {
		leaf.LogRequest.DurationSeconds = int(expr.Window.Duration(0) / 1000)
	}
	return leaf, nil
}
