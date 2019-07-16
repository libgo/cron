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

type ScheduleMode int

const (
	minutely ScheduleMode = iota
	hourly
	daily
	weekly
	monthly
	yearly
)

type Schedule struct {
	mode   ScheduleMode
	second int
	minute int
	hour   int
	day    int
	month  int
	dow    int
}

func (s *Schedule) Next() time.Time {
	now := time.Now()

	switch s.mode {
	case yearly:
		r := time.Date(now.Year(), time.Month(s.month), s.day, s.hour, s.minute, s.second, 0, now.Location())
		if r.Before(now) {
			r = r.AddDate(1, 0, 0)
		}
		return r
	case monthly:
		r := time.Date(now.Year(), now.Month(), s.day, s.hour, s.minute, s.second, 0, now.Location())
		if r.Before(now) {
			r = r.AddDate(0, 1, 0)
		}
		return r
	case weekly:
		r := time.Date(now.Year(), now.Month(), now.Day(), s.hour, s.minute, s.second, 0, now.Location())
		if r.Before(now) {
			r = r.AddDate(0, 0, 1)
		}
		for {
			if r.Weekday() == time.Weekday(s.dow) {
				break
			}
			r = r.AddDate(0, 0, 1)
		}
		return r
	case daily:
		r := time.Date(now.Year(), now.Month(), now.Day(), s.hour, s.minute, s.second, 0, now.Location())
		if r.Before(now) {
			r = r.AddDate(0, 0, 1)
		}
		return r
	case hourly:
		r := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), s.minute, s.second, 0, now.Location())
		if r.Before(now) {
			r = r.Add(1 * time.Hour)
		}
		return r
	case minutely:
		r := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), s.second, 0, now.Location())
		if r.Before(now) {
			r = r.Add(1 * time.Minute)
		}
		return r
	}
	return now // never come here
}

// i is input timing string, format should be "second minute hour day month dow"
// minutely: 2 * * * * *     "*:*:02"
// hourly: 2 10 * * * *      "*:10:02"
// daily: 2 10 3 * * *       "3:10:02"
// monthly: 2 10 3 1 * *     "1st 3:10:02"
// yearly: 2 10 3 1 2 *      "Feb 1st 3:10:02"
// weekly: 2 10 3 * * 2      "Tue 3:10:02"
func Parse(i string) (*Schedule, error) {
	s := strings.Fields(i)
	if len(s) != 6 {
		return nil, ErrFormat
	}

	second, err := strconv.Atoi(s[0])
	if err != nil || second < 0 || second > 59 {
		return nil, ErrFormat
	}

	// minutely
	if s[1] == "*" { // minutely
		if s[2] != "*" || s[3] != "*" || s[4] != "*" || s[5] != "*" {
			return nil, ErrFormat
		}
		return &Schedule{minutely, second, 0, 0, 0, 0, 0}, nil
	}

	minute, err := strconv.Atoi(s[1])
	if err != nil || minute < 0 || minute > 59 {
		return nil, ErrFormat
	}

	// hourly
	if s[2] == "*" {
		if s[3] != "*" || s[4] != "*" || s[5] != "*" {
			return nil, ErrFormat
		}
		return &Schedule{hourly, second, minute, 0, 0, 0, 0}, nil
	}

	hour, err := strconv.Atoi(s[2])
	if err != nil || hour < 0 || hour > 23 {
		return nil, ErrFormat
	}

	// error check
	if s[3] == "*" && s[4] != "*" {
		return nil, ErrFormat
	}
	if s[3] != "*" && s[5] != "*" {
		return nil, ErrFormat
	}
	if s[4] != "*" && s[5] != "*" {
		return nil, ErrFormat
	}

	if s[3] == "*" {
		// daily
		if s[5] == "*" {
			return &Schedule{daily, second, minute, hour, 0, 0, 0}, nil
		}

		// weekly
		dow, err := strconv.Atoi(s[5])
		if err != nil || dow < 0 || dow > 6 {
			return nil, ErrFormat
		}
		return &Schedule{weekly, second, minute, hour, 0, 0, dow}, nil
	}

	day, err := strconv.Atoi(s[3])
	if err != nil || day < 1 || day > 31 {
		return nil, ErrFormat
	}

	// monthly
	if s[monthly] == "*" {
		return &Schedule{monthly, second, minute, hour, day, 0, 0}, nil
	}

	// yearly
	month, err := strconv.Atoi(s[4])
	if err != nil || month < 1 || month > 12 {
		return nil, ErrFormat
	}
	return &Schedule{yearly, second, minute, hour, day, month, 0}, nil
}
