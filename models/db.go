package models

import (
    "database/sql"
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
