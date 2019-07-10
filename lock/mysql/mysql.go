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

func conf(uri string) mysqlx.Conf {
	return mysqlx.Conf{
		DSN:             uri,
		MaxOpenConns:    16,
		MaxIdleConns:    8,
		ConnMaxLifetime: time.Minute * 15,
	}
}

// Open a new mysql locker instance
func (m *Mysql) Open(uri string) lock.Locker {
	db := mysqlx.Register("db", mysqlx.Conf{
		DSN:             uri,
		MaxOpenConns:    16,
		MaxIdleConns:    8,
		ConnMaxLifetime: time.Minute * 15,
	})

	return &Mysql{db: db}
}

// Lock with GET_LOCK
func (m *Mysql) Lock(n string) error {
	query := fmt.Sprintf(`SELECT GET_LOCK("job_lock_%s", 2)`, n)
	var success bool

	if err := m.db.QueryRow(query).Scan(&success); err != nil {
		return fmt.Errorf("lock job %s error: %s", n, err.Error())
	}

	if !success {
		return fmt.Errorf("lock job %s error", n)
	}

	return nil
}

// Unlock with RELEASE_LOCK
func (m *Mysql) Unlock(n string) error {
	query := fmt.Sprintf(`SELECT RELEASE_LOCK("job_lock_%s")`, n)
	if _, err := m.db.Exec(query); err != nil {
		return fmt.Errorf("unlock job %s error: %s", n, err.Error())
	}
	return nil
}
