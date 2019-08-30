package mysql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/libgo/cron/lock"
	"github.com/libgo/logx"
	"github.com/libgo/mysqlx"
)

func init() {
	lock.Register("mysql", &Mysql{})
}

type Mysql struct {
	db *sqlx.DB
}

// Open a new mysql locker instance
func (m *Mysql) Open(uri string) lock.Locker {
	db := mysqlx.Register("db", mysqlx.Conf{
		DSN:             uri,
		MaxOpenConns:    2,
		MaxIdleConns:    1,
		ConnMaxLifetime: time.Minute * 5,
	})

	return &Mysql{db: db}
}

// Lock with GET_LOCK
func (m *Mysql) Lock(n string) error {
	query := fmt.Sprintf(`SELECT GET_LOCK("job_lock_%s", 2)`, n)
	var success bool

	if err := m.db.QueryRow(query).Scan(&success); err != nil {
		logx.Errorf("cron: lock job %s error at mysql db: %s", n, err.Error())
		return lock.ErrLock
	}

	if !success {
		return lock.ErrLock
	}

	return nil
}

// Unlock with RELEASE_LOCK
func (m *Mysql) Unlock(n string) error {
	query := fmt.Sprintf(`SELECT RELEASE_LOCK("job_lock_%s")`, n)
	if _, err := m.db.Exec(query); err != nil {
		logx.Errorf("cron: unlock job %s error at mysql db: %s", n, err.Error())
		return lock.ErrUnlock
	}

	return nil
}
