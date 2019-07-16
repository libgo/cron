package null

import (
	"github.com/libgo/cron/lock"
)

func init() {
	lock.Register("null", &Null{})
}

// Null is a locker that not supporting unique executing of job
type Null struct{}

func (n *Null) Open(uri string) lock.Locker {
	return n
}

func (n *Null) Lock(ns string) error {
	return nil
}

func (n *Null) Unlock(ns string) error {
	return nil
}
