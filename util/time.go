package util

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/spf13/cast"
)

var durationRE = regexp.MustCompile("^(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?$")

func GetDurationFromAge(durationStr string) (time.Duration, error) {
	switch durationStr {
	case "0":
		return 0, nil
	case "":
		return 0, errors.New("empty duration string")
	}
	matches := durationRE.FindStringSubmatch(durationStr)
	if matches == nil {
		return 0, fmt.Errorf("not a valid duration string: %q", durationStr)
	}
	var dur time.Duration
	var overflowErr error
	m := func(pos int, mult time.Duration) {
		if matches[pos] == "" {
			return
		}
		n, _ := strconv.Atoi(matches[pos])
		if n > int((1<<63-1)/mult/time.Millisecond) {
			overflowErr = errors.New("duration out of range")
		}
		d := time.Duration(n) * time.Millisecond
		dur += d * mult
		if dur < 0 {
			overflowErr = errors.New("duration out of range")
		}
	}
	m(2, 1000*60*60*24*365) // y
	m(4, 1000*60*60*24*7)   // w
	m(6, 1000*60*60*24)     // d
	m(8, 1000*60*60)        // h
	m(10, 1000*60)          // m
	m(12, 1000)             // s
	m(14, 1)                // ms
	return dur, overflowErr
}

func FloatStringToTime(timeFloat string) time.Time {
	sec, dec := math.Modf(cast.ToFloat64(timeFloat))
	return time.Unix(int64(sec), int64(dec*(1e9)))
}
