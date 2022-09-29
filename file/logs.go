package file

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/util"
	"github.com/spf13/cast"
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
	NodePattern string
	Nodes       []string
	Process     string
}

type PodSearchParams struct {
	NamespacePattern string
	Namespaces       []string
	Pod              string
	Container        string
}

type LogSearch struct {
	LogType       string // audit | event | pod | node
	TargetPattern string // audit | event | <namespace_pattern> | <node_pattern>
	Targets       []string
	// AuditSearchParams AuditSearchParams
	EventSearchParams EventSearchParams
	PodSearchParams   PodSearchParams
	NodeSearchParams  NodeSearchParams
	Keywords          []string
	DurationSeconds   int
	EndTime           time.Time
	IsCounting        bool
}

type Result struct {
	IsCounting bool
	Logs       []string
	Count      int
}

func GetLogs(logSearch LogSearch) (Result, error) {
	fmt.Println("logSearch=", logSearch)
	var command string
	var out string
	var err error

	if logSearch.LogType == "pod" || logSearch.LogType == "node" {
		if logSearch.TargetPattern == "" {
			logSearch.TargetPattern = "*"
		}
		logSearch.TargetPattern = strings.ReplaceAll(logSearch.TargetPattern, ".*", "*")
	}

	// list namespace directories
	command = fmt.Sprintf("ls -d %s/%s/%s", config.GetLogRoot(), logSearch.LogType, logSearch.TargetPattern)
	out, _, err = util.Execute(command)
	if err != nil {
		// return Result{}, fmt.Errorf("not found matched namespace directories: err=%s", err)
		if logSearch.IsCounting {
			return Result{IsCounting: true, Count: 0}, nil
		}
		return Result{IsCounting: false, Logs: []string{}}, nil
	}
	dirs := strings.Split(strings.TrimRight(out, "\n"), "\n")

	// list files
	command = fmt.Sprintf("ls %s", strings.Join(funk.Map(dirs, func(x string) string { return x + "/*.log" }).([]string), " "))
	// fmt.Println("getLogs: command=", command)
	out, code, err := util.Execute(command)
	if code != 2 && err != nil {
		return Result{}, fmt.Errorf("not found matched files: err=%s", err)
	}
	files := strings.Split(strings.TrimRight(out, "\n"), "\n")
	// fmt.Println("getLogs: files=", files)
	logFileSearchs := getLogFilesSearchs(files, logSearch.DurationSeconds, logSearch.EndTime)
	// fmt.Println("getLogs: logFileSearchs=", logFileSearchs)
	result, err := getLogsFromLogFileSearchs(logFileSearchs, logSearch)
	if err != nil {
		return Result{}, err
	}
	// fmt.Println("logs=", logs)
	return result, nil
}

func getLogFilesSearchs(filepaths []string, durationSeconds int, endTime time.Time) []logFileSearch {
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
	// fmt.Printf("getLogFilesSearchs: files=%#v\n", files)
	for _, filepath := range filepaths {
		parts := strings.Split(filepath, "/")
		files = append(files, parts[len(parts)-1])
	}

	// unique files
	logFileSearchMap := map[string]logFileSearch{}
	for _, file := range files {
		fmt.Println("file=", file)
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

	// reverse sort by file
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

func getLogsFromLogFileSearchs(logFileSearchs []logFileSearch, logSearch LogSearch) (Result, error) {
	logs := []string{}
	count := 0
	limit := config.GetLimit()
	for _, logFileSearch := range logFileSearchs {
		newResult, err := getLogsFromLogFileSearch(logFileSearch, logSearch, limit)
		if err != nil {
			log.Printf("warning on getPodLogsFrompodLogFileSearch: err=%s\n", err)
			break
		}
		if logSearch.IsCounting {
			count += newResult.Count
			limit = config.GetLimit() - count
		} else {
			logs = append(newResult.Logs, logs...)
			limit = config.GetLimit() - len(logs)
		}
		if limit < 1 {
			break
		}
	}
	if logSearch.IsCounting {
		return Result{IsCounting: true, Count: count}, nil
	}
	return Result{Logs: logs}, nil
}

func getLogsFromLogFileSearch(logFileSearch logFileSearch, logSearch LogSearch, limit int) (Result, error) {
	// fmt.Println("logFileSearch=", logFileSearch)
	files := ""
	for _, target := range logFileSearch.Targets {
		files += fmt.Sprintf("%s/%s/%s/%s ", config.GetLogRoot(), logSearch.LogType, target, logFileSearch.File)
	}
	sort := ""
	// if len(logFileSearch.Targets) > 1 {
	// 	sort = "| sort"
	// }
	egrepTimePatterns := ""
	if len(logFileSearch.TimePatterns) > 0 {
		egrepTimePatterns = fmt.Sprintf("| egrep ^'%s'", strings.Join(logFileSearch.TimePatterns, "|"))
	}
	grepLabels := ""
	switch logSearch.LogType {
	case "audit":
	case "event":
		namespace := logSearch.EventSearchParams.Namespace
		typ := logSearch.EventSearchParams.Type
		reason := logSearch.EventSearchParams.Reason
		object := logSearch.EventSearchParams.Object
		count := logSearch.EventSearchParams.Count

		if namespace == "" && typ == "" && reason == "" && object == "" && count == "" {
			break // no grep
		}

		if namespace == "" {
			namespace = ".*"
		}
		if typ == "" {
			typ = ".*"
		}
		if reason == "" {
			reason = ".*"
		}
		if object == "" {
			object = ".*"
		}
		if count == "" {
			count = ".*"
		}
		grepLabels = fmt.Sprintf(`| grep 'Z\[%s|%s|%s|%s|%s\]'`, namespace, typ, reason, object, count)

	case "node":
		process := logSearch.NodeSearchParams.Process
		if process != "" {
			grepLabels = fmt.Sprintf(`| grep 'Z\[.*|%s\]'`, process)
		}
	case "pod":
		pod := logSearch.PodSearchParams.Pod
		container := logSearch.PodSearchParams.Container
		if pod == "" && container == "" {
			break // no break
		}
		if pod == "" {
			pod = ".*"
		}
		if container == "" {
			container = ".*"
		}
		grepLabels = fmt.Sprintf(`| grep 'Z\[.*|%s|%s\]'`, pod, container)

	}
	grepLabels = strings.ReplaceAll(grepLabels, `.`, `[^|]`)
	// grepLabels = strings.ReplaceAll(grepLabels, `*`, `\*`)
	grepLabels = strings.ReplaceAll(grepLabels, `+`, `\+`)
	grepKeyword := ""
	for _, keyword := range logSearch.Keywords {
		if len(keyword) < 1 {
			return Result{}, errors.New("empty keyword")
		}
		grepKeyword = fmt.Sprintf(" | grep %s", "'"+strings.ReplaceAll(keyword, "'", "'\"'\"'")+"'")
	}
	countCommand := ""
	if logSearch.IsCounting {
		countCommand = "| wc -l"
	}
	command := fmt.Sprintf("timeout 30 cat %s %s %s %s %s | tail -%d %s", files, sort, egrepTimePatterns, grepLabels, grepKeyword, limit, countCommand)
	// fmt.Println("command=", command)
	out, code, err := util.Execute(command)
	if code != 1 && err != nil {
		return Result{}, err
	}
	if code == 1 || out == "" {
		if logSearch.IsCounting {
			return Result{IsCounting: logSearch.IsCounting}, nil
		}
		return Result{Logs: []string{}}, nil
	}
	if logSearch.IsCounting {
		return Result{IsCounting: true, Count: cast.ToInt(strings.TrimSpace(out))}, nil
	}
	return Result{Logs: strings.Split(strings.TrimRight(out, "\n"), "\n")}, nil
}
