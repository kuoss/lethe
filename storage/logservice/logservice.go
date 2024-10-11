package logservice

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/storage/fileservice"
	"github.com/kuoss/lethe/storage/logservice/logmodel"
	"github.com/kuoss/lethe/storage/logservice/match"
)

type LogService struct {
	config      *config.Config
	fileService *fileservice.FileService
}

func New(cfg *config.Config, fileService *fileservice.FileService) *LogService {
	return &LogService{cfg, fileService}
}

func (s *LogService) SelectLog(sel *model.LogSelector) (log model.Log, warnings model.Warnings, err error) {
	// type
	if sel.Name != "node" && sel.Name != "pod" {
		return log, warnings, fmt.Errorf("unknown logType: %s", sel.Name)
	}

	// targets
	targets, err := s.getTargets(sel, &warnings)
	if err != nil {
		return log, warnings, fmt.Errorf("getTargets err: %w", err)
	}

	// files
	files := s.getFiles(sel, targets, &warnings)

	// mfs
	mfs, err := match.GetMatchFuncSet(sel)
	if err != nil {
		return log, warnings, fmt.Errorf("getMatchFuncSet err: %w", err)
	}

	// log
	log = s.getLogFromFiles(sel, files, mfs, &warnings)
	return log, warnings, nil
}

// targets
func (s *LogService) getTargets(sel *model.LogSelector, warnings *model.Warnings) ([]string, error) {
	all, err := s.fileService.List(sel.Name)
	if err != nil {
		return nil, fmt.Errorf("list err: %w", err)
	}
	matcher, ws, err := match.GetTargetMatcher(sel)
	if err != nil {
		return nil, fmt.Errorf("target matcher err: %w", err)
	}
	*warnings = append(*warnings, ws...)

	targets := []string{}
	for _, t := range all {
		if matcher.Matches(filepath.Base(t)) {
			targets = append(targets, t)
		}
	}
	return targets, nil
}

func (s *LogService) getFiles(sel *model.LogSelector, targets []string, ws *model.Warnings) []string {
	var files []string
	for _, target := range targets {
		list, err := s.fileService.List(target)
		if err != nil {
			*ws = append(*ws, fmt.Errorf("list err: %w", err))
			continue
		}
		for _, item := range list {
			if isFileInTimeRange(item, &sel.TimeRange) {
				files = append(files, item)
			}
		}
	}
	return files
}

func isFileInTimeRange(file string, tr *model.TimeRange) bool {
	name := filepath.Base(file)
	fileStart, err := time.Parse(time.RFC3339, strings.Replace(name[0:13], "_", "T", 1)+":00:00Z")
	if err != nil {
		return false
	}
	fileEnd := fileStart.Add(time.Hour) // per hour for one logs
	return tr.Start.Before(fileEnd) && tr.End.After(fileStart)
}

// log
func (s *LogService) getLogFromFiles(sel *model.LogSelector, files []string, mfs *match.MatchFuncSet, warnings *model.Warnings) model.Log {
	limit := s.config.Query.Limit
	logLines := []model.LogLine{}
	sort.Strings(files)
	for _, file := range files {
		s.addLogLinesFromFile(sel, &logLines, file, mfs, warnings)
		if len(logLines) > limit {
			break
		}
	}
	return model.Log{
		Name:  sel.Name,
		Lines: logLines,
	}
}

func (s *LogService) addLogLinesFromFile(sel *model.LogSelector, logLines *[]model.LogLine, file string, mfs *match.MatchFuncSet, warnings *model.Warnings) {
	limit := s.config.Query.Limit
	sc, err := s.fileService.Scanner(file)
	if err != nil {
		*warnings = append(*warnings, err)
		return
	}
	for sc.Scan() {
		addLogLine(sel, logLines, sc.Text(), mfs, warnings)
		if len(*logLines) > limit {
			return
		}
	}
}

func addLogLine(sel *model.LogSelector, logLines *[]model.LogLine, line string, mfs *match.MatchFuncSet, warnings *model.Warnings) {
	pos := strings.IndexRune(line, '[')
	if pos < 0 {
		*warnings = append(*warnings, fmt.Errorf("no time separator"))
		return
	}
	tim := line[:pos]
	parsedTime, err := time.Parse(time.RFC3339Nano, tim)
	if err != nil {
		*warnings = append(*warnings, fmt.Errorf("time parse err: %w", err))
		return
	}
	if sel.TimeRange.Start.After(parsedTime) || sel.TimeRange.End.Before(parsedTime) {
		return // skip
	}
	rest := line[pos+1:]

	// labels]log
	pos = strings.IndexRune(rest, ']')
	if pos < 0 {
		*warnings = append(*warnings, fmt.Errorf("no log separator"))
		return
	}
	labels := strings.Split(rest[:pos], "|")
	log := rest[pos+2:]

	// label match
	if len(labels) != len(mfs.LabelMatchFuncs)+1 {
		*warnings = append(*warnings, fmt.Errorf("label count mismatch"))
		return
	}
	for i, f := range mfs.LabelMatchFuncs {
		if f == nil {
			continue
		}
		if !f(labels[i+1]) {
			return
		}
	}

	// line match
	for _, f := range mfs.LineMatchFuncs {
		if !f(log) {
			return
		}
	}

	switch sel.Name {
	case "node":
		*logLines = append(*logLines, logmodel.NodeLog{Time: tim, Node: labels[0], Process: labels[1], Log: log})
		return
	case "pod":
		*logLines = append(*logLines, logmodel.PodLog{Time: tim, Namespace: labels[0], Pod: labels[1], Container: labels[2], Log: log})
		return
	}
	*warnings = append(*warnings, fmt.Errorf("addLogLine: unknown log type: %s", sel.Name))
}
