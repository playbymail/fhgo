// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package sqlite3

import (
	"database/sql"
	_ "embed"
	"errors"
	"github.com/playbymail/fhgo/domains"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

var (
	//go:embed schema.sql
	schemaDDL string
)

// DatabaseCreate creates a new database.
// Returns an error if the database already exists.
func DatabaseCreate(path string, force bool) error {
	sb, err := os.Stat(path)
	// if the stat fails because the file doesn't exist, we're okay.
	// if it fails for any other reason, it's an error.
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("[sqlite3] %q: %s\n", path, err)
		return err
	}
	// it is an error if the path exists and is not a regular file.
	if sb != nil && (sb.IsDir() || !sb.Mode().IsRegular()) {
		log.Printf("[sqldb] %q: is a folder\n", path)
		return domains.ErrInvalidPath
	}
	// it is an error if the database already exists unless force is true.
	// in that case, we remove the database so that we can create it again.
	if sb != nil { // database file exists
		if !force {
			// we're not forcing the creation of a new database so this is an error
			return domains.ErrDatabaseExists
		}
		log.Printf("[sqlite3] removing %s\n", path)
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	// create the database
	log.Printf("[sqlite3] creating %s\n", path)
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Printf("[sqlite3] %s\n", err)
		return err
	}
	defer db.Close()

	// confirm that the database driver supports foreign keys
	checkPragma := "PRAGMA" + " foreign_keys = ON"
	if rslt, err := db.Exec(checkPragma); err != nil {
		log.Printf("[sqlite3] error: foreign keys are disabled\n")
		return domains.ErrForeignKeysDisabled
	} else if rslt == nil {
		log.Printf("[sqlite3] error: foreign keys pragma failed\n")
		return domains.ErrPragmaReturnedNil
	}

	// create the schema
	if _, err := db.Exec(schemaDDL); err != nil {
		log.Printf("[sqlite3] failed to initialize schema\n")
		log.Printf("[sqlite3] %v\n", err)
		return errors.Join(domains.ErrCreateSchema, err)
	}

	log.Printf("[sqlite3] created %s\n", path)

	return nil
}
