package models

import (
    "fmt"
    "encoding/json"
    "crypto/rand"
)

//Product struct preserves single DB row of product
type Product struct{
    Serial string       `json:"serial"`
    Name string         `json:"name"`
    Company string      `json:"company"`
    Place int           `json:"place"`
    Column int          `json:"column"`
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

//AllProducts handles SQL request to get all products from DB
func (db *DB) AllProducts() (Serializable, error) {
    rows, err := db.Query("select * from Goods")
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
            fmt.Println(err)
            return nil, err
        }
        products = append(products, p)
    }
    return products, nil
}
//InsertProduct adds a new product entry in DB
func (db *DB) InsertProduct(parse []byte) error {
    var p Product
    err := json.Unmarshal(parse, &p)
    if (err != nil) {
        return err
    }
    req, err := db.Prepare("insert into Goods values($1,$2,$3,$4,$5,$6)")
    defer req.Close()
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
    res, err := req.Exec(
        p.Serial,
        p.Name,
        p.Company,
        p.Row,
        p.Column,
        p.Place,
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
//DeleteProduct deletes product with selected ID from DB
func (db *DB) DeleteProduct(parse []byte) error {
    var p Product
    err := json.Unmarshal(parse, &p)
    if (err != nil) {
        return err
    }
    id := p.Serial
    req, err := db.Prepare("delete from Goods where Serial = ?")
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
