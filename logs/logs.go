package logs

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/kuoss/lethe/storage/driver"
	"github.com/kuoss/lethe/storage/driver/factory"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/thoas/go-funk"
)

// NEW
type logFileSearch struct {
	File         string
	Targets      []string
	TimePatterns []string
}

// type AuditSearchParams struct {
// }

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

// Todo now Patternable can handle only one '*' and starts with substring
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
	if ps.Patterned() {
		//todo
		//nginx-*  vs nginx-.* (regex)
		if string(ps)[strings.IndexRune(string(ps), '*')-1] == '.' {
			return string(ps)[0 : strings.IndexRune(string(ps), '*')-1]
		}
		return string(ps)[0:strings.IndexRune(string(ps), '*')]
	}
	return string(ps)
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
	Filter            Filter
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

func dirExist(dirs []string, dir string) bool {
	for _, dirFullPath := range dirs {
		parts := strings.Split(dirFullPath, "/")
		if parts[len(parts)-1] == dir {
			return true
		}
	}
	return false
}

func (ls *LogStore) GetLogs(logSearch LogSearch) (Result, error) {
	//fmt.Printf("logSearch= %+v", logSearch)
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
		now := config.GetNow()
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

	sort.Sort(sort.StringSlice(files))
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
				if search.Filter.match(line) {
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
				if search.Filter.match(line) {
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

func matchedTimePatternFiles(filepaths []string, durationSeconds int, endTime time.Time) []logFileSearch {
	// durationSeconds & timeRange => EndTimeAgo & endTimeAgo
	if endTime.IsZero() {
		endTime = config.GetNow()
	}
	var startTime time.Time
	if durationSeconds == 0 {
		startTime = endTime.Add(time.Duration(-100*24) * time.Hour)
	} else {
		startTime = endTime.Add(time.Duration(-durationSeconds) * time.Second)
	}
	// fmt.Println("\nstartTime=", startTime, "endTime=", endTime)
	// all files
	files := []string{}
	// fmt.Printf("matchedTimePatternFiles: files=%#v\n", files)
	for _, filepath := range filepaths {
		parts := strings.Split(filepath, "/")
		files = append(files, parts[len(parts)-1])
	}

	// unique files
	logFileSearchMap := map[string]logFileSearch{}
	for _, file := range files {
		//fmt.Println("logs=", logs)
		// RFC3339 "2006-01-02T15:04:05Z07:00"
		fileStartTime, err := time.Parse(time.RFC3339, strings.Replace(file[0:13], "_", "T", 1)+":00:00Z")
		if err != nil {
			log.Println(err)
			continue
		}
		fileEndTime := fileStartTime.Add(time.Duration(3599) * time.Second)
		// fmt.Println("fileStartTime=", fileStartTime, "fileEndTime=", fileEndTime)
		if startTime.After(fileEndTime) || endTime.Before(fileStartTime) {
			// fmt.Println("out of range")
			continue
		}
		if startTime.Before(fileStartTime) && endTime.After(fileEndTime) {
			// fmt.Println("in range")
			logFileSearchMap[file] = logFileSearch{File: file, Targets: []string{}, TimePatterns: []string{}}
			continue
		}
		inFileStartTime := fileStartTime
		if inFileStartTime.Before(startTime) {
			inFileStartTime = startTime
		}
		inFileEndTime := fileEndTime
		if inFileEndTime.After(endTime) {
			inFileEndTime = endTime
		}
		// fmt.Println("inFileStartTime=", inFileStartTime, "inFileEndTime=", inFileEndTime)
		timePatterns := []string{}
		for t := inFileStartTime; t.Equal(inFileEndTime) || t.Before(inFileEndTime); t = t.Add(time.Duration(1) * time.Second) {
			// fmt.Println("timePattern=", strings.Replace(t.UTC().String()[0:19], " ", "T", 1))
			timePatterns = append(timePatterns, strings.Replace(t.UTC().String()[0:19], " ", "T", 1))
		}
		// fmt.Println("timePatterns=", timePatterns)
		logFileSearchMap[file] = logFileSearch{File: file, Targets: []string{}, TimePatterns: reduceTimePatterns(timePatterns)}
	}
	if len(filepaths) > 0 {
		for _, filepath := range filepaths {
			parts := strings.Split(filepath, "/")
			file := parts[len(parts)-1]
			target := parts[len(parts)-2]
			if entry, ok := logFileSearchMap[file]; ok {
				entry.Targets = append(entry.Targets, target)
				logFileSearchMap[file] = entry
			}
		}
	}
	// fmt.Printf("logFileSearchMap=%#v\n", logFileSearchMap)

	// reverse sort by logs
	keys := funk.Keys(logFileSearchMap).([]string)
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	logFileSearchList := []logFileSearch{}
	for _, key := range keys {
		logFileSearchList = append(logFileSearchList, logFileSearchMap[key])
	}
	// fmt.Printf("logFileSearchList=%#v\n", logFileSearchList)
	return logFileSearchList
}

func reduceTimePatterns(timePatterns []string) []string {
	if len(timePatterns) == 3600 {
		return []string{} // 1 hours (3600 seconds) doesn't need timePatterns
	}
	timePatterns = reduceTimePatternsGroup(timePatterns, 15, 600) // 12:3x:xx (600 seconds)
	timePatterns = reduceTimePatternsGroup(timePatterns, 16, 60)  // 12:34:xx (60 seconds)
	timePatterns = reduceTimePatternsGroup(timePatterns, 18, 10)  // 12:34:5x (10 seconds)
	return timePatterns
}

func reduceTimePatternsGroup(timePatterns []string, length, full int) []string {
	temps := funk.UniqString(funk.Map(timePatterns, func(x string) string {
		if len(x) < length {
			return ""
		}
		return x[:length]
	}).([]string))
	for _, temp := range temps {
		if temp == "" {
			continue
		}
		idxs := []int{}
		for idx, timePattern := range timePatterns {
			if len(timePattern) >= length && temp == timePattern[:length] {
				idxs = append(idxs, idx)
			}
		}
		if len(idxs) != full {
			continue
		}
		timePatterns = append(timePatterns[:idxs[0]], timePatterns[idxs[full-1]+1:]...)
		timePatterns = append(timePatterns, temp)
	}
	return timePatterns
}

// todo
// go to PodLog type private fucntion
func parseHierarchyPod(line string) (namespace, pod, container string, err error) {
	pos := line[strings.IndexRune(line, '[')+1 : strings.IndexRune(line, ']')]
	parts := strings.Split(pos, "|")
	if len(parts) != 3 {
		return namespace, pod, container, errors.New(fmt.Sprintf("log line not follow [{ns}|{pod}|{container}]. : %s", line))
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
		return node, process, errors.New(fmt.Sprintf("log line not follow [{node}|{process}]. : %s", line))
	}
	node, process = parts[0], parts[1]
	return node, process, nil
}
