package mysql

import (
	"github.com/jetuuuu/hl_homework/database"
)

type DB struct {
	db *database.DB
}

func New(db *database.DB) *DB {
	return &DB{db}
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Underlying() *database.DB {
	return db.db
}
