package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrFormat = fmt.Errorf("input schedule string format error")
)

type Schedule struct {
	second int
	minute int
	hour   int
	day    int
	month  int
	dow    int
}

func (s *Schedule) Next() time.Time {
	now := time.Now()

	if s.month != 0 {
		r := time.Date(now.Year(), time.Month(s.month), s.day, s.hour, s.minute, s.second, 0, now.Location())
		if time.Month(s.month) > now.Month() {
			return r
		}
		return r.AddDate(1, 0, 0)
	}

	if s.day != 0 {
		r := time.Date(now.Year(), now.Month(), s.day, s.hour, s.minute, s.second, 0, now.Location())
		if s.day > now.Day() {
			return r
		}
		return r.AddDate(0, 1, 0)
	}

	r := time.Date(now.Year(), now.Month(), now.Day(), s.hour, s.minute, s.second, 0, now.Location())
	if r.Before(now) {
		r = r.AddDate(0, 0, 1)
	}

	if s.dow != 0 {
		dow := s.dow
		if dow == 7 { // convert to time.Weekday format for Sunday
			dow = 0
		}

		for {
			if r.Weekday() != time.Weekday(dow) {
				r.AddDate(0, 0, 1)
			} else {
				return r
			}
		}
	}

	return r
}

// i is input string, format should be "second minute hour day month dow"
// second = [0..59]
// minute = [0..59]
// hour = [0..23]
// day = [1..31]   0 means ignore
// month = [1..12]   0 means ignore
// dow = [1..7]   0 means ignore, 7 is sunday
// only support 6 models
// minutely: 2 0 0 0 0 0     "x:x:02"
// hourly: 2 10 0 0 0 0      "x:10:02"
// daily: 2 10 3 0 0 0       "3:10:02"
// monthly: 2 10 3 2 0 0     "1st 3:10:02"
// yearly: 2 10 3 2 1 0      "Jan 1st 3:10:02"
// weekly: 2 10 3 0 0 2      "Tue 3:10:02"
func Parse(i string) (*Schedule, error) {
	i = strings.Replace(i, "*", "-1", -1)
	s := strings.Split(i, " ")
	if len(s) != 6 {
		return nil, ErrFormat
	}

	second, err := strconv.Atoi(s[0])
	if err != nil || second < 0 || second > 59 {
		return nil, ErrFormat
	}

	minute, err := strconv.Atoi(s[1])
	if err != nil || minute < 0 || minute > 59 {
		return nil, ErrFormat
	}

	hour, err := strconv.Atoi(s[2])
	if err != nil || hour < 0 || hour > 23 {
		return nil, ErrFormat
	}

	day, err := strconv.Atoi(s[3])
	if err != nil || day < 0 || day > 31 {
		return nil, ErrFormat
	}

	month, err := strconv.Atoi(s[4])
	if err != nil || month < 0 || month > 31 {
		return nil, ErrFormat
	}

	dow, err := strconv.Atoi(s[5])
	if err != nil || dow < 0 || dow > 31 {
		return nil, ErrFormat
	}

	if dow != 0 && day != 0 {
		return nil, ErrFormat
	}

	if dow != 0 && month != 0 {
		return nil, ErrFormat
	}

	return &Schedule{second, minute, hour, day, month, dow}, nil
}
