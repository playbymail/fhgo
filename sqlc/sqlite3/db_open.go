// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package sqlite3

import (
	"context"
	"database/sql"
	"github.com/playbymail/fhgo/domains"
	"log"
	"os"
)

type DatabaseCloser func() error

// DatabaseOpen opens an existing database.
// Returns an error if the path is not a directory, or if the database does not exist.
// Caller must use the DatabaseCloser to close the database connection when done.
func DatabaseOpen(path string, ctx context.Context) (*Queries, DatabaseCloser, error) {
	// it is an error if the database does not already exist and is not a file.
	sb, err := os.Stat(path)
	if err != nil {
		log.Printf("[sqlite3] %q: %s\n", path, err)
		return nil, func() error { return nil }, err
	} else if sb.IsDir() || !sb.Mode().IsRegular() {
		log.Printf("[sqlite3] %q: %s\n", path, err)
		return nil, func() error { return nil }, domains.ErrInvalidPath
	}
	log.Printf("[sqlite3] opening %s\n", path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	// confirm that the database driver supports foreign keys
	checkPragma := "PRAGMA" + " foreign_keys = ON"
	if rslt, err := db.Exec(checkPragma); err != nil {
		_ = db.Close()
		log.Printf("[sqlite3] error: foreign keys are disabled\n")
		return nil, func() error { return nil }, domains.ErrForeignKeysDisabled
	} else if rslt == nil {
		_ = db.Close()
		log.Printf("[sqlite3] error: foreign keys pragma failed\n")
		return nil, func() error { return nil }, domains.ErrPragmaReturnedNil
	}

	// return a Query which wraps the database handle. this means the caller
	// can't directly access the database handle. we must return a function
	// that will close the database handle for them.
	return New(db), func() error {
		var err error
		if db != nil {
			err, db = db.Close(), nil
		}
		return err
	}, nil
}
