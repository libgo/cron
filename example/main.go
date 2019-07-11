package main

import (
	"flag"
	"time"

	"github.com/libgo/cron"
	_ "github.com/libgo/cron/lock/mysql"
	"github.com/libgo/logx"
)

type PrintJob struct {
	s string
}

func (j *PrintJob) Name() string {
	return "printing"
}

func (j *PrintJob) Run() {
	logx.Infof("printing %s", j.s)
	time.Sleep(10 * time.Second)
}

var (
	s string
)

func init() {
	flag.StringVar(&s, "d", "dummy string", "dummy string")
}

func main() {
	flag.Parse()

	err := cron.Locker("mysql", "root:passWORD@tcp(192.168.10.191:3306)/dolphin")
	if err != nil {
		logx.Errorf("init locker error: %s", err.Error())
	}

	cron.Add("0 24 10 0 0 0", &PrintJob{s: s})
	cron.Add("0 25 10 0 0 0", &PrintJob{s: s})
	cron.Add("0 26 10 0 0 0", &PrintJob{s: s})
	cron.Add("0 27 10 0 0 0", &PrintJob{s: s})
	cron.Add("0 28 10 0 0 0", &PrintJob{s: s})
	cron.Add("0 29 10 0 0 0", &PrintJob{s: s})
	cron.Add("0 30 10 0 0 0", &PrintJob{s: s})
	cron.Run()

	c := make(chan bool)
	<-c
}
