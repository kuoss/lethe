package queryservice

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/lethe/clock"
	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/letheql"
	"github.com/kuoss/lethe/letheql/model"
	"github.com/kuoss/lethe/storage/logservice"
)

type QueryService struct {
	engine *letheql.Engine
}

func New(cfg *config.Config, logService *logservice.LogService) *QueryService {
	return &QueryService{
		engine: letheql.NewEngine(cfg, logService),
	}
}

func (s *QueryService) Query(ctx context.Context, qs string, tr model.TimeRange) *letheql.Result {
	if reflect.ValueOf(tr).IsZero() {
		now := clock.Now()
		tr = model.TimeRange{
			Start: now.Add(-1 * time.Minute),
			End:   now,
		}
	}

	qs = toLetheQL(qs)
	qry, err := s.engine.NewRangeQuery(qs, tr.Start, tr.End)
	if err != nil {
		return &letheql.Result{Err: err}
	}
	res := qry.Exec(ctx)
	if res.Err != nil {
		logger.Errorf("exec err: %s", res.Err.Error())
	}
	return res
}

func toLetheQL(input string) string {
	if !strings.HasPrefix(input, "{") {
		return input
	}

	re := regexp.MustCompile(`^\{([^}]+)\}(.*)$`)
	matches := re.FindStringSubmatch(input)

	if len(matches) < 3 {
		return input
	}

	contentInBraces := matches[1] // e.g., `job="pod",namespace="namespace01"``
	remainderString := matches[2] // e.g., ` |= "hello world"`

	// Regex to find and capture the 'job="VALUE"' part.
	// This regex is made more robust to handle spaces around commas or at the start/end.
	reType := regexp.MustCompile(`(?:^|\s*,\s*)job="([^"]+)"(?:$|\s*,\s*)`)
	jobMatches := reType.FindStringSubmatch(contentInBraces)

	// If no 'job' attribute is found, return the original input as per requirements.
	if len(jobMatches) < 2 { // Changed from 3 to 2 as the new regex has only one capturing group for the value
		return input
	}

	extractedType := jobMatches[1] // The captured VALUE of job (e.g., "pod", "node")

	// Split by commas (with optional spaces) to process individual parts
	parts := regexp.MustCompile(`\s*,\s*`).Split(contentInBraces, -1)

	filteredParts := []string{}
	jobRemoved := false
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if !jobRemoved && strings.HasPrefix(trimmedPart, "job=") {
			jobRemoved = true // Mark that 'job' has been processed/removed
			continue
		}
		if trimmedPart != "" {
			filteredParts = append(filteredParts, trimmedPart)
		}
	}

	tempContent := strings.TrimSpace(strings.Join(filteredParts, ","))
	if tempContent == "" {
		return extractedType + remainderString
	} else {
		return fmt.Sprintf("%s{%s}%s", extractedType, tempContent, remainderString)
	}
}
