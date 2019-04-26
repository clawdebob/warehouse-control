package models

import (
    "database/sql"
    "fmt"
    "reflect"
    "strings"
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
    AllOrders() (Serializable, error)
    InsertProduct([]byte) error
    DeleteProduct([]byte) error
    InsertPerson([]byte) error
    DeletePerson([]byte) error
    InsertOrder([]byte) error
    FilterProduct([]byte, string) (Serializable, error)
    FilterPerson([]byte, string) (Serializable, error)
    EditProduct([]byte) error
    EditPerson([]byte) error
}
//DB describes struct that implements Datastore
type DB struct {
    *sql.DB
    filter func(string, interface{}) (string, error)
    update func(string, interface{}) (string, error)
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
    return &DB{db,
            queryWrapper("SELECT * FROM %s WHERE","%s LIKE '%s%%'", "AND"),
            queryWrapper("UPDATE %s SET","%s = '%s'", ","),
         }, nil
}

func queryWrapper (sQuery string, sCondition string, sep string) (func (string, interface{}) (string, error)) {
    return func (table string, e interface{}) (string, error) {
        finalQuery := fmt.Sprintf(sQuery, table)
        query := make([]string, 0)
        typ := reflect.TypeOf(e)
        val := reflect.ValueOf(e)
        for c := 0; c < typ.NumField(); c++ {
            field := typ.Field(c)
            value := val.Field(c)
            sqlName, ok := field.Tag.Lookup("sql")
            if !ok {
                sqlName = field.Name
            }
            switch field.Type.Kind() {
                case reflect.Int:
                    if value.Int() != 0 {
                        query = append(query, fmt.Sprintf(" %s = %d", sqlName, value.Int()))
                    }
                    break;
                case reflect.String:
                    if value.String() != "" {
                        query = append(query, fmt.Sprintf(" " + sCondition, sqlName, value.String()))
                    }
                    break;
                default:
                    return "", fmt.Errorf("unhandled type")
            }
        }
        finalQuery += strings.Join(query, " " + sep)
        return finalQuery, nil
    }
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
