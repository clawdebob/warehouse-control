package models

import (
    "fmt"
    "encoding/json"
    "crypto/rand"
)

//Product represents single DB row of product
type Product struct{
    Serial string       `json:"serial"`
    Name string         `json:"name"`
    Company string      `json:"company" sql:"Manufacturer"`
    Place int           `json:"place"`
    Column int          `json:"column" sql:"Columm"`
    Row int             `json:"row"`
}
//FromJSON for Product converts JSON to Product entry
func (p *Product) FromJSON(parse []byte) (error) {
    err := json.Unmarshal(parse, p)
    return err
}

//Products is a type that describes Products array
type Products []Product

//ToJSON for Products converts Product array to JSON
func (p Products) ToJSON() (string, error) {
    json, err := json.MarshalIndent(p, "", "    ")
    return string(json), err
}

func (db *DB) getProductsQuery(q string) (Serializable, error) {
    fmt.Println(q)
    rows, err := db.Query(q)
    if err != nil {
        return nil, err
    }

    defer rows.Close()
    products := Products{}

    for rows.Next(){
        p:= Product {}
        err:= rows.Scan(
            &p.Serial,
            &p.Name,
            &p.Company,
            &p.Row,
            &p.Column,
            &p.Place,
         )
        if err != nil {
            return nil, err
        }
        products = append(products, p)
    }
    return products, nil
}

//AllProducts handles SQL request to get all products from DB
func (db *DB) AllProducts() (Serializable, error) {
    return db.getProductsQuery("SELECT * FROM Goods")
}
//InsertProduct adds a new product entry in DB
func (db *DB) InsertProduct(parse []byte) error {
    var p Product
    err := json.Unmarshal(parse, &p)
    if (err != nil) {
        return err
    }
    b := make([]byte, 16)
    _, err = rand.Read(b)
    if err != nil {
        return err
    }
    serial := fmt.Sprintf("%x", b[0:4])
    p.Serial = serial
    return db.execEntity(
        "INSERT INTO Goods VALUES($1,$2,$3,$4,$5,$6)",
        p.Serial,
        p.Name,
        p.Company,
        p.Row,
        p.Column,
        p.Place,
    )
}
//EditProduct edits specified parameters of one product entry
func (db *DB) EditProduct(data []byte) error {
    var p Product
    err := json.Unmarshal(data, &p)
    if err != nil {
        return err
    }
    finalQuery, err := db.update("Goods", p)
    if err != nil {
        return err
    }
    if p.Serial != "" {
        finalQuery += fmt.Sprintf(" WHERE Serial = '%s'", p.Serial)
    } else {
        return fmt.Errorf("serial is empty")
    }

    return db.execEntity(finalQuery)
}
//DeleteProduct deletes product with selected Serial from DB
func (db *DB) DeleteProduct(parse []byte) error {
    var p Product
    err := json.Unmarshal(parse, &p)
    if err != nil {
        return err
    }
    id := p.Serial
    return db.execEntity("DELETE FROM Goods WHERE Serial = ?", id)
}
//FilterProduct filters Good's rows according to specified filters
func (db *DB) FilterProduct(data []byte, sort string) (Serializable, error) {
    var p Product
    sortBy := " ORDER BY Name"
    err := json.Unmarshal(data, &p)
    if err != nil {
        return nil, err
    }
    finalQuery, err := db.filter("Goods", p)
    if err != nil {
        return nil, err
    }
    if sort == "desc" {
        sortBy += " DESC"
    }
    finalQuery += sortBy
    return db.getProductsQuery(finalQuery)
}
