package models

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3" //SQLite driver
)

//NewDB creates connection to specified DB
func NewDB(databaseName string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", databaseName)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
