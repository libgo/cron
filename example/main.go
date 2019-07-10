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

	err := cron.SetLocker("mysql", "root:ddg1208@tcp(192.168.10.191:3306)/dolphin")
	if err != nil {
		logx.Errorf("set locker error: %s", err.Error())
	}

	cron.Add("0 20 21 0 0 0", &PrintJob{s: s})
	cron.Add("0 21 21 0 0 0", &PrintJob{s: s})
	cron.Add("0 22 21 0 0 0", &PrintJob{s: s})
	cron.Add("0 23 21 0 0 0", &PrintJob{s: s})
	// cron.Add("0 36 20 0 0 4", &PrintJob{s: s})
	// cron.Add("0 36 20 12 0 0", &PrintJob{s: s})
	// cron.Add("0 36 20 12 9 0", &PrintJob{s: s})
	// cron.Add("0 36 20 12 5 0", &PrintJob{s: "tesing runing"})
	cron.Run()

	c := make(chan bool)
	<-c
}
