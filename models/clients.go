package models

import (
    "encoding/json"
    "strings"
    "fmt"
)
//Person struct represents single DB row of client
type Person struct {
    ID int              `json:"id" sql:"Id"`
    Name string         `json:"name"`
    Type int            `json:"type" sql:"ClientType"`
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

func (db *DB) getPersonsQuery(q string) (Serializable, error) {
    fmt.Println(q)
    rows, err := db.Query(q)
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

//AllPersons handles SQL request to get all persons from DB
func (db *DB) AllPersons() (Serializable, error) {
    return db.getPersonsQuery("SELECT * FROM Clients")
}

//InsertPerson adds a new person entry in DB
func (db *DB) InsertPerson(parse []byte) error {
    var p Person
    err := json.Unmarshal(parse, &p)
    if (err != nil) {
        return err
    }
    return db.execEntity(
        "INSERT INTO Clients(Name,ClientType,Address,Email,PhoneNumber) VALUES($1,$2,$3,$4,$5)",
        p.Name,
        p.Type,
        p.Address,
        p.Email,
        p.PhoneNumber,
    )
}
//DeletePerson deletes person with selected ID from DB
func (db *DB) DeletePerson(parse []byte) error {
    var p Person
    err := json.Unmarshal(parse, &p)
    if (err != nil) {
        return err
    }
    return db.execEntity("DELETE FROM Clients WHERE Id = ?", p.ID)
}
//EditPerson edits peson with selected ID
func (db *DB) EditPerson(data []byte) error{
    var p Person
    finalQuery := "UPDATE Clients SET"
    query := make([]string, 0)
    err := json.Unmarshal(data, &p)
    if (err != nil) {
        return err
    }
    if p.Name != "" {
        query = append(query, fmt.Sprintf(" Name = '%s'", p.Name))
    }
    if p.Address != "" {
        query = append(query, fmt.Sprintf(" Address = '%s'", p.Address))
    }
    if p.Email != "" {
        query = append(query, fmt.Sprintf(" Email = '%s'", p.Email))
    }
    if p.PhoneNumber != "" {
        query = append(query, fmt.Sprintf(" PhoneNumber = '%s'", p.PhoneNumber))
    }
    if p.Type != 0 {
        query = append(query, fmt.Sprintf(" ClientType = %d", p.Type))
    }
    finalQuery += strings.Join(query, " ,")
    if p.ID != 0 {
        finalQuery += fmt.Sprintf(" WHERE Id = %d", p.ID)
    } else {
        return fmt.Errorf("id is invalid")
    }

    return db.execEntity(finalQuery)
}

//FilterPerson filters Client's rows according to specified filters
func (db *DB) FilterPerson(data []byte, sort string) (Serializable, error) {
    var p Person
    sortBy := " ORDER BY Name"
    err := json.Unmarshal(data, &p)
    if (err != nil) {
        return nil, err
    }
    finalQuery, err := db.constructFilterQuery("Clients", p)
    if (err != nil) {
        return nil, err
    }
    if sort == "desc" {
        sortBy += " DESC"
    }
    finalQuery += sortBy
    return db.getPersonsQuery(finalQuery)
}
