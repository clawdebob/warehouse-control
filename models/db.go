package models

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3" //SQLite driver
)

//Serializable describes DB entities than can be serialized to json
type Serializable interface {
    ToJSON() (string, error)
}
//Datastore contains meta data about all methods that our db should implement
type Datastore interface {
    AllProducts() (Serializable, error)
    AllPersons() (Serializable, error)
    InsertProduct([]byte) error
    DeleteProduct([]byte) error
    InsertPerson([]byte) error
    DeletePerson([]byte) error
    FilterProduct([]byte, string) (Serializable, error)
    EditProduct([]byte) error
}
//DB describes struct that implements Datastore
type DB struct {
    *sql.DB
}

//NewDB creates connection to specified DB
func NewDB(databaseName string) (*DB, error) {
    db, err := sql.Open("sqlite3", databaseName)
    db.SetMaxOpenConns(1)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return &DB{db}, nil
}

func (db *DB) execEntity(q string, args ...interface{}) error {
    req, err := db.Prepare(q)
    defer req.Close()
    if (err != nil) {
        return err
    }
    res, err := req.Exec( args...)
    if (err != nil) {
        return err
    }
    rc, err := res.RowsAffected()
    if (err != nil) {
        return err
    }
    if (rc == 0) {
        return fmt.Errorf("warning!!! 0 rows affected")
    }
    return nil
}
