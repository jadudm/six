package dyndb

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type DynDB struct {
	DB       *sql.DB
	Path     string
	Filename string
}

// I want to use the DB with
func NewDynDB(path string, filename string) (*DynDB, error) {
	os.MkdirAll(path, 0744)
	db, err := sql.Open("sqlite", filepath.Join(path, filename))
	if err != nil {
		return nil, errors.New(fmt.Sprint(err))
	}
	dyndb := DynDB{
		DB:       db,
		Path:     path,
		Filename: filename,
	}
	return &dyndb, nil
}
