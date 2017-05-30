package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/acoshift/session"
)

// Errors
var (
	ErrDBRequired    = errors.New("sql: db required")
	ErrTableRequired = errors.New("sql: table required")
)

// Config is the sql store config
type Config struct {
	DB              *sql.DB
	Table           string
	CleanupInterval time.Duration
}

// New creates new sql store
func New(config Config) session.Store {
	db := config.DB
	table := config.Table
	cleanupInterval := config.CleanupInterval

	if db == nil {
		panic(ErrDBRequired)
	}
	if len(table) == 0 {
		panic(ErrTableRequired)
	}

	db.Exec(fmt.Sprintf(`
		create table if not exists %s (
			k text,
			v blob,
			e timestamp,
			primary key (k),
			index (e)
		);
	`, table))
	// ignore create table error

	s := &sqlStore{
		db:              db,
		cleanupInterval: cleanupInterval,
		getQuery:        fmt.Sprintf(`select v, e, now() from %s where k = $1`, table),
		setQuery:        fmt.Sprintf(`insert into %s (k, v, e) values ($1, $2, $3) on conflict (k) do update set v = excluded.v, k = excluded.k`, table),
		delQuery:        fmt.Sprintf(`delete from %s where k = $1`, table),
		expQuery:        fmt.Sprintf(`update %s set e = $2 where k = $1`, table),
		delExpiredQuery: fmt.Sprintf(`delete from %s where e <= now()`, table),
	}
	if cleanupInterval > 0 {
		go s.cleanupWorker()
	}
	return s
}

type sqlStore struct {
	db              *sql.DB
	cleanupInterval time.Duration
	getQuery        string
	setQuery        string
	delQuery        string
	expQuery        string
	delExpiredQuery string
}

var errNotFound = errors.New("sql: session not found")

func (s *sqlStore) cleanupWorker() {
	// add small delay before start worker
	time.Sleep(5 * time.Second)
	for {
		s.db.Exec(s.delExpiredQuery)
		time.Sleep(s.cleanupInterval)
	}
}

func (s *sqlStore) Get(key string) ([]byte, error) {
	var bs []byte
	var exp *time.Time
	var now time.Time
	err := s.db.QueryRow(s.getQuery, key).Scan(&bs, &exp, &now)
	if err != nil {
		return nil, err
	}
	if exp != nil && exp.Before(now) {
		s.Del(key)
		return nil, errNotFound
	}
	return bs, nil
}

func (s *sqlStore) Set(key string, value []byte, ttl time.Duration) error {
	var exp *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		exp = &t
	}
	_, err := s.db.Exec(s.setQuery, key, value, exp)
	return err
}

func (s *sqlStore) Del(key string) error {
	_, err := s.db.Exec(s.delQuery, key)
	return err
}

func (s *sqlStore) Exp(key string, ttl time.Duration) error {
	var exp *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		exp = &t
	}
	_, err := s.db.Exec(s.expQuery, key, exp)
	return err
}
