package models

import (
    "encoding/json"
    "fmt"
)
//Person struct represents single DB row of client
type Person struct {
    ID int              `json:"id"`
    Name string         `json:"name"`
    Type int            `json:"type"`
    Address string      `json:"address"`
    Email string        `json:"email"`
    PhoneNumber string  `json:"phone"`
}
//Persons type represents Person array
type Persons []Person

//ToJSON for Products converts Product array to JSON
func (p Persons) ToJSON() (string, error) {
    json, err := json.MarshalIndent(p, "", "    ")
    return string(json), err
}

//AllPersons handles SQL request to get all persons from DB
func (db *DB) AllPersons() (Serializable, error) {
    rows, err := db.Query("select * from Clients")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    persons := Persons{}

    for rows.Next(){
        p:= Person {}
        err:= rows.Scan(
            &p.ID,
            &p.Name,
            &p.Type,
            &p.Address,
            &p.Email,
            &p.PhoneNumber,
         )
        if err != nil {
            return nil, err
        }
        persons = append(persons, p)
    }
    return persons, nil
}

//InsertPerson adds a new person entry in DB
func (db *DB) InsertPerson(parse []byte) error {
    var p Person
    err := json.Unmarshal(parse, &p)
    if (err != nil) {
        return err
    }
    req, err := db.Prepare("insert into Clients(Name,ClientType,Address,Email,PhoneNumber) values($1,$2,$3,$4,$5)")
    defer req.Close()
    if (err != nil) {
        return err
    }
    res, err := req.Exec(
        p.Name,
        p.Type,
        p.Address,
        p.Email,
        p.PhoneNumber,
    )
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
//DeletePerson deletes person with selected ID from DB
func (db *DB) DeletePerson(parse []byte) error {
    var p Person
    err := json.Unmarshal(parse, &p)
    if (err != nil) {
        return err
    }
    id := p.ID
    req, err := db.Prepare("delete from Clients where Id = ?")
    defer req.Close()
    if (err != nil) {
        return err
    }
    res, err := req.Exec(id)
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
