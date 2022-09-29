package letheql

import (
	"errors"
	"regexp"
	"strings"
)

var partSeparator = regexp.MustCompile("[{}]")

func ParseQuery(queryString string) (ParsedQuery, error) {
	queryString = strings.TrimSpace(queryString)
	var query ParsedQuery
	parts := partSeparator.Split(queryString, -1)
	partsLength := len(parts)
	if partsLength != 1 && partsLength != 3 {
		return query, errors.New("invalid query")
	}
	switch strings.TrimSpace(parts[0]) {
	case "pod":
		query.Type = "pod"
	case "node":
		query.Type = "node"
	default:
		return query, errors.New("invalid query type")
	}
	// fmt.Printf("parts=%#v\n", parts)
	// fmt.Printf("partsLength=%#v\n", partsLength)
	if partsLength == 1 {
		return query, nil
	}
	labels, err := ParseLabels(strings.TrimSpace(parts[1]))
	if err != nil {
		// fmt.Println("err=", err)
		return query, errors.New("invalid label format")
	}
	query.Labels = labels
	query.Keyword = strings.TrimSpace(parts[2])
	// fmt.Printf("query=%#v\n", query)
	return query, nil
}

func ParseLabels(labelsString string) ([]Label, error) {
	var labels []Label
	labelsString = strings.TrimSpace(labelsString)
	if labelsString == "" {
		return labels, nil
	}
	parts := strings.Split(labelsString, ",")
	for _, part := range parts {
		subparts := strings.Split(part, "=")
		if len(subparts) == 0 {
			continue
		}
		if len(subparts) != 2 {
			return labels, errors.New("invalid label format")
		}

		key := subparts[0]
		key = strings.ReplaceAll(key, `"`, ``)
		key = strings.TrimSpace(key)

		value := subparts[1]
		value = strings.ReplaceAll(value, `"`, ``)
		value = strings.TrimSpace(value)

		labels = append(labels, Label{
			Key:   key,
			Value: value,
		})
	}
	return labels, nil
}
