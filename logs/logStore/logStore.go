package logStore

import (
	"bufio"
	"fmt"
	"github.com/kuoss/lethe/logs/filter"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/storage/driver"
	"github.com/kuoss/lethe/storage/driver/factory"

	"github.com/kuoss/lethe/config"
)

type EventSearchParams struct {
	Namespace string
	Type      string
	Reason    string
	Object    string
	Count     string
}

type NodeSearchParams struct {
	Nodes   []string
	Node    PatternedString
	Process PatternedString
}

type PodSearchParams struct {
	Namespaces []string
	Namespace  PatternedString
	Pod        PatternedString
	Container  PatternedString
}

// now Patternable can handle only one '*' and starts with substring
// exampel) nignx-*, namespace*,
type Patternable interface {
	Patterned() bool
}
type PatternedString string

func (ps PatternedString) Patterned() bool {
	return strings.Contains(string(ps), "*")
}

func (ps PatternedString) PatternMatch(s string) bool {
	if string(ps) == "" {
		return true
	}
	if ps.Patterned() {
		return strings.Contains(s, ps.withoutPattern())
	}
	return ps.withoutPattern() == s
}

// withoutPattern return pattern removed string,
// if it has not pattern return its own string
func (ps PatternedString) withoutPattern() string {
	psString := string(ps)
	if !ps.Patterned() {
		return psString
	}
	//todo
	//nginx-*  vs nginx-.* (regex)
	pos := strings.IndexRune(psString, '*')
	if pos > 0 && psString[pos-1] == '.' {
		return psString[0 : pos-1]
	}
	return psString[0:pos]
}

type LogSearch struct {
	LogType       LogType         // audit | event | pod | node
	TargetPattern PatternedString // audit | event | <namespace_pattern> | <node_pattern>
	Targets       []string
	// AuditSearchParams AuditSearchParams
	EventSearchParams EventSearchParams
	PodSearchParams   PodSearchParams
	NodeSearchParams  NodeSearchParams
	Keyword           string
	DurationSeconds   int
	EndTime           time.Time
	StartTime         time.Time
	IsCounting        bool
	Filter            filter.Filter
}

type Result struct {
	IsCounting bool
	Logs       []string
	Count      int
}

type LogStore struct {
	driver driver.StorageDriver
}

func New() *LogStore {
	d, _ := factory.Get("filesystem", map[string]interface{}{"RootDirectory": config.GetLogRoot()})
	return &LogStore{driver: d}
}

func (ls *LogStore) GetLogs(logSearch LogSearch) (Result, error) {
	// fmt.Printf("logSearch= %+v", logSearch)
	//fmt.Printf("root directory  from driver: %s\n", ls.driver.RootDirectory())

	logTypePath := filepath.Join(ls.driver.RootDirectory(), logSearch.LogType.GetName())
	targets, err := ls.driver.List(logTypePath)
	//fmt.Printf("log Type: %s, targets: %v", logTypePath, targets)

	if err != nil {
		return Result{IsCounting: false, Logs: []string{}}, nil
	}

	var matchedTarget []string

	if logSearch.TargetPattern.Patterned() {
		var patternMatched []string
		for _, t := range targets {
			//todo
			_, candidates := filepath.Split(t)
			if strings.Contains(candidates, logSearch.TargetPattern.withoutPattern()) {
				patternMatched = append(patternMatched, candidates)
			}
		}
		matchedTarget = patternMatched
	} else {
		matchedTarget = append(matchedTarget, string(logSearch.TargetPattern))
	}

	rangeParamInit(&logSearch)

	//fmt.Printf("log Type: %s, matched targets: %v", logTypePath, matchedTarget)
	// from here only check matchedTarget
	var timeFilteredFiles []string
	for _, dir := range matchedTarget {
		timeFilteredFiles = append(timeFilteredFiles, timeFilter(filepath.Join(logTypePath, dir), &logSearch, ls.driver)...)
	}

	//fmt.Printf("log Type: %s, timfiltered targets: %v\n", logTypePath, timeFilteredFiles)
	logs := logFromTarget(timeFilteredFiles, logSearch, config.GetLimit(), ls.driver)

	sort.SliceStable(*logs, func(i, j int) bool {
		l := *logs
		x, _ := timeParse(l[i])
		y, _ := timeParse(l[j])
		return x.Before(y)
	})

	return Result{Logs: *logs, Count: len(*logs)}, nil
}

func rangeParamInit(search *LogSearch) {
	if search.EndTime.IsZero() {
		now := clock.GetNow()
		search.EndTime = now
	}

	if search.DurationSeconds == 0 {
		search.StartTime = search.EndTime.Add(time.Duration(-100*24) * time.Hour) // 10 days ago
	} else {
		search.StartTime = search.EndTime.Add(time.Duration(-search.DurationSeconds) * time.Second)
	}
}

// todo
// limit check here?
func logFromTarget(files []string, search LogSearch, limit int, driver driver.StorageDriver) (logs *[]string) {

	logs = &[]string{}

	sort.Strings(files)
	for _, file := range files {
		limit -= checkTarget(file, search, logs, driver)
		if limit < 0 {
			return logs
		}
	}
	return logs
}

// This function process line by line
// very Dependent on performance
func checkTarget(file string, search LogSearch, logs *[]string, driver driver.StorageDriver) (logSize int) {
	switch search.LogType.GetName() {
	case "audit":
	case "event":
	case NODE_TYPE:

		//todo process()
		rc, _ := driver.Reader(file)
		sc := bufio.NewScanner(rc)
		for sc.Scan() {
			// 2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]
			line := sc.Text()
			//todo error handling
			timeFromLog, _ := timeParse(line)
			if search.StartTime.After(timeFromLog) || search.EndTime.Before(timeFromLog) {
				continue
			}

			node, process, err := parseHierarchyNode(line)
			if err != nil {
				return 0
			}

			if !search.NodeSearchParams.Node.PatternMatch(node) || !search.NodeSearchParams.Process.PatternMatch(process) {
				continue
			}

			//todo filtering here?
			if search.Filter != nil {
				if search.Filter.Match(line) {
					*logs = append(*logs, line)
				}
			} else {
				*logs = append(*logs, line)
			}
		}
		return len(*logs)

	case POD_TYPE:
		//todo process()
		rc, _ := driver.Reader(file)
		sc := bufio.NewScanner(rc)
		for sc.Scan() {
			// 2009-11-10T23:00:00.000000Z[namespace01|nginx-deployment-75675f5897-7ci7o|nginx] hello world [ddd12wewe]
			line := sc.Text()
			//todo error handling
			timeFromLog, _ := timeParse(line)
			if search.StartTime.After(timeFromLog) || search.EndTime.Before(timeFromLog) {
				continue
			}

			ns, pod, container, err := parseHierarchyPod(line)
			if err != nil {
				return 0
			}

			if !search.PodSearchParams.Namespace.PatternMatch(ns) || !search.PodSearchParams.Pod.PatternMatch(pod) || !search.PodSearchParams.Container.PatternMatch(container) {
				continue
			}
			//todo filtering here?
			if search.Filter != nil {
				if search.Filter.Match(line) {
					*logs = append(*logs, line)
				}
			} else {
				*logs = append(*logs, line)
			}
		}
		return len(*logs)
	}
	return 0
}

func timeParse(line string) (time.Time, error) {
	s := line[0:strings.IndexRune(line, '[')]
	parsed, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}

func timeFilter(directory string, search *LogSearch, driver driver.StorageDriver) (filteredFiles []string) {

	files, err := driver.List(directory)
	if err != nil {
		return
	}

	for _, file := range files {
		filename := filepath.Base(file)
		if rangeCheckFromFilename(filename, search.StartTime, search.EndTime) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles
}

func rangeCheckFromFilename(name string, start time.Time, end time.Time) bool {
	fileStart, err := time.Parse(time.RFC3339, strings.Replace(name[0:13], "_", "T", 1)+":00:00Z")
	if err != nil {
		//to do
		return false
	}

	fileEnd := fileStart.Add(time.Duration(3599) * time.Second) // per hour for one logs

	if start.After(fileEnd) || end.Before(fileStart) {
		// out of range
		return false
	}
	// except the whole range return true
	return true
}

// todo
// go to PodLog type private fucntion
func parseHierarchyPod(line string) (namespace, pod, container string, err error) {
	pos := line[strings.IndexRune(line, '[')+1 : strings.IndexRune(line, ']')]
	parts := strings.Split(pos, "|")
	if len(parts) != 3 {
		return namespace, pod, container, fmt.Errorf("log line not follow [{ns}|{pod}|{container}]. : %s", line)
	}
	namespace, pod, container = parts[0], parts[1], parts[2]
	return namespace, pod, container, nil
}

// todo
// go to NodeLog type private fucntion
func parseHierarchyNode(line string) (node, process string, err error) {
	meta := line[strings.IndexRune(line, '[')+1 : strings.IndexRune(line, ']')]
	parts := strings.Split(meta, "|")
	if len(parts) != 2 {
		return node, process, fmt.Errorf("log line not follow [{node}|{process}]. : %s", line)
	}
	node, process = parts[0], parts[1]
	return node, process, nil
}
