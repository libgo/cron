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
	// Run must > 2 second
	Run()
}

type entry struct {
	schedule Schedule
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

// SetLock config lock target
// n is the name of lock service: "mysql", "etcd", "redis"
// d is the connect string to the service
func SetLocker(n string, d string) error {
	l, err := lock.Open(n, d)
	if err != nil {
		return err
	}

	locker = l
	return nil
}

// Add job to crontab
func Add(i string, j Job) error {
	if running {
		return ErrRunning
	}

	s, err := Parse(i)
	if err != nil {
		return err
	}

	entries = append(entries, entry{schedule: *s, job: j})
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
			for {
				select {
				case <-tm.C:
					err := locker.Lock(j.Name())
					if err == nil {
						run(j)
						locker.Unlock(j.Name())
					}

					tm.Reset(s.Next().Sub(time.Now()))
				}
			}
		}()
	}
	return nil
}

// run with recovery
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
