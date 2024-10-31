// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package sqlc implements a data store using the Sqlite3 database.

package sqlc

//go:generate sqlc generate

import (
	"context"
	_ "embed"
	"github.com/playbymail/fhgo/domains"
	"github.com/playbymail/fhgo/sqlc/sqlite3"
	"github.com/playbymail/fhgo/stdfs"
)

type DB struct {
	path   string
	ctx    context.Context
	q      *sqlite3.Queries
	Closer sqlite3.DatabaseCloser
}

func Open(path string, ctx context.Context) (db *DB, err error) {
	if exists, err := stdfs.IsFileExists(path); err != nil {
		return nil, err
	} else if !exists {
		return nil, domains.ErrDatabaseMissing
	}
	db = &DB{path: path, ctx: ctx}
	db.q, db.Closer, err = sqlite3.DatabaseOpen(path, db.ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) Close() error {
	return db.Closer()
}
