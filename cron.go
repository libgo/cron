package cron

import (
	"fmt"
	"runtime"
	"time"

	"github.com/libgo/cron/lock"
	"github.com/libgo/logx"
)

// Job
type Job interface {
	Name() string
	Run() // For job, the run routine should > 2s, add a time.Sleep(2 * time.Second) is a good idea.
}

type entry struct {
	schedule *Schedule
	job      Job
}

var (
	ErrRunning = fmt.Errorf("cron has already been running, cannot do anything more")
	ErrLock    = fmt.Errorf("cron lock is not nil")
)

var (
	entries             = []entry{}
	running             = false
	locker  lock.Locker = nil
)

// Locker init lock container
// n is the name of lock service: "mysql", "etcd", "redis"
// uri is the connect string to the service
func Locker(n string, uri string) error {
	l, err := lock.Open(n, uri)
	if err != nil {
		return err
	}

	locker = l
	return nil
}

// Add job to crontab
// i is input timing string, format should be "second minute hour day month dow"
// minutely: 2 * * * * *     "*:*:02"
// hourly: 2 10 * * * *      "*:10:02"
// daily: 2 10 3 * * *       "3:10:02"
// monthly: 2 10 3 1 * *     "1st 3:10:02"
// yearly: 2 10 3 1 2 *      "Feb 1st 3:10:02"
// weekly: 2 10 3 * * 2      "Tue 3:10:02"
func Add(i string, j Job) error {
	if running {
		return ErrRunning
	}

	s, err := Parse(i)
	if err != nil {
		return err
	}

	entries = append(entries, entry{schedule: s, job: j})
	return nil
}

// Run all jobs in crontab
func Run() error {
	if running {
		return ErrRunning
	}

	if locker == nil {
		return ErrLock
	}

	for _, e := range entries {
		s := e.schedule
		j := e.job

		go func() {
			tm := time.NewTimer(s.Next().Sub(time.Now()))
			logx.Infof("cron: schedule job '%s' at %s", j.Name(), s.Next().String())
			for {
				select {
				case <-tm.C:
					err := locker.Lock(j.Name())
					if err == nil {
						run(j)
						locker.Unlock(j.Name())
					}

					tm.Reset(s.Next().Sub(time.Now()))
					logx.Infof("cron: schedule job '%s' at %s", j.Name(), s.Next().String())
				}
			}
		}()
	}
	return nil
}

// run with recovery. For cron job, should not break the service.
func run(j Job) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			logx.Errorf("cron: panic running job: %v\n%s", r, buf)
		}
	}()

	j.Run()
}
