package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/jetuuuu/hl_homework/log"
)

type DB struct {
	db         *sql.DB
	tx         *sql.Tx
	instanceID string
}

func Open(driverName, dbinfo, instanceID string) (_ *DB, err error) {
	db, err := sql.Open(driverName, dbinfo)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return New(db, instanceID), nil
}

func New(db *sql.DB, instanceID string) *DB {
	return &DB{db: db, instanceID: instanceID}
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) InTransaction() bool {
	return db.tx != nil
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	defer logQuery(ctx, query, args, db.instanceID)(&err)

	if db.tx != nil {
		return db.tx.ExecContext(ctx, query, args...)
	}
	return db.db.ExecContext(ctx, query, args...)
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (_ *sql.Rows, err error) {
	defer logQuery(ctx, query, args, db.instanceID)(&err)

	if db.tx != nil {
		return db.tx.QueryContext(ctx, query, args...)
	}

	return db.db.QueryContext(ctx, query, args...)
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	defer logQuery(ctx, query, args, db.instanceID)(nil)

	if db.tx != nil {
		return db.tx.QueryRowContext(ctx, query, args...)
	}

	return db.db.QueryRowContext(ctx, query, args...)
}

func (db *DB) Transact(ctx context.Context, iso sql.IsolationLevel, txFunc func(*DB) error) (err error) {
	opts := &sql.TxOptions{Isolation: iso}
	return db.transact(ctx, opts, txFunc)
}

func (db *DB) transact(ctx context.Context, opts *sql.TxOptions, txFunc func(*DB) error) (err error) {
	if db.InTransaction() {
		return errors.New("a DB Transact function was called on a DB already in a transaction")
	}
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("db.BeginTx(): %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			if txErr := tx.Commit(); txErr != nil {
				err = fmt.Errorf("tx.Commit(): %w", txErr)
			}
		}
	}()

	dbtx := New(db.db, db.instanceID)
	dbtx.tx = tx

	if err := txFunc(dbtx); err != nil {
		return fmt.Errorf("txFunc(tx): %w", err)
	}
	return nil
}

var QueryLoggingDisabled bool

var queryCounter int64

type queryEndLogEntry struct {
	ID              string
	Query           string
	Args            string
	DurationSeconds float64
	Error           string `json:",omitempty"`
}

func logQuery(ctx context.Context, query string, args []interface{}, instanceID string) func(*error) {
	if QueryLoggingDisabled {
		return func(*error) {}
	}
	const maxlen = 300

	var r []rune
	for _, c := range query {
		if c == '\n' {
			c = ' '
		}
		if len(r) == 0 || !unicode.IsSpace(r[len(r)-1]) || !unicode.IsSpace(c) {
			r = append(r, c)
		}
	}
	query = string(r)
	if len(query) > maxlen {
		query = query[:maxlen] + "..."
	}

	uid := generateLoggingID(instanceID)

	const (
		maxArgs   = 20
		maxArgLen = 50
	)
	var argStrings []string
	for i := 0; i < len(args) && i < maxArgs; i++ {
		s := fmt.Sprint(args[i])
		if len(s) > maxArgLen {
			s = s[:maxArgLen] + "..."
		}
		argStrings = append(argStrings, s)
	}
	if len(args) > maxArgs {
		argStrings = append(argStrings, "...")
	}
	argString := strings.Join(argStrings, ", ")

	log.Debugf(ctx, "%s %s args=%s", uid, query, argString)
	start := time.Now()
	return func(errp *error) {
		dur := time.Since(start)
		if errp == nil {
			log.Debugf(ctx, "%s done", uid)
		} else {
			entry := queryEndLogEntry{
				ID:              uid,
				Query:           query,
				Args:            argString,
				DurationSeconds: dur.Seconds(),
			}

			if *errp == nil {
				log.Debug(ctx, "entry", entry)
			} else {
				entry.Error = (*errp).Error()
				log.Error(ctx, "entry", entry)
			}
		}
	}
}

func generateLoggingID(instanceID string) string {
	if instanceID == "" {
		instanceID = "local"
	} else {
		instanceID = instanceID[len(instanceID)-4:]
	}
	n := atomic.AddInt64(&queryCounter, 1)
	return fmt.Sprintf("%s-%d", instanceID, n)
}
