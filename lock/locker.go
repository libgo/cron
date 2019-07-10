package lock

import "fmt"

var (
	ErrNoLocker = fmt.Errorf("invalid locker")
	lockers     = make(map[string]Locker)
)

// Locker interface
type Locker interface {
	Open(string) Locker
	Lock(string) error
	Unlock(string) error
}

// Open a new locker with name and uri to mysql/etcd/redis
func Open(n string, uri string) (Locker, error) {
	l, ok := lockers[n]
	if !ok {
		return nil, ErrNoLocker
	}

	return l.Open(uri), nil
}

// Register mysql/etcd/redis locker
func Register(n string, l Locker) {
	lockers[n] = l
}
