package main

import (
	"flag"
	"os"
	"time"

	"github.com/libgo/cron"
	_ "github.com/libgo/cron/lock/mysql"
	_ "github.com/libgo/cron/lock/null"
	_ "github.com/libgo/cron/lock/redis"
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
	time.Sleep(3 * time.Second)
}

var (
	i string
	s string
	l string
	u string
)

func init() {
	flag.StringVar(&i, "i", "0 * * * * *", "timing spec")
	flag.StringVar(&s, "s", "dummy", "namespace")
	flag.StringVar(&l, "l", "mysql", "locker, now support mysql, redis, null")
	flag.StringVar(&u, "u", "root:passWORD@tcp(192.168.10.191:3306)/dolphin", `uri to the locker service, now support redis and mysql.
mysql pattern should like: 'root:passWORD@tcp(192.168.10.191:3306)/dolphin'
redis pattern should like: 'redis://:passWORD@192.168.10.145:6379'`)
}

func main() {
	flag.Parse()

	err := cron.Locker(l, u)
	if err != nil {
		logx.Errorf("init locker error: %s", err.Error())
		os.Exit(1)
	}

	err = cron.Add(i, &PrintJob{s: s + " job"})
	if err != nil {
		logx.Errorf("add cron job error: %s", err.Error())
		os.Exit(1)
	}

	cron.Run()

	c := make(chan bool)
	<-c
}
