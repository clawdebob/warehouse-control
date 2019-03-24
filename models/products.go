package models

import (
    "fmt"
    "encoding/json"
)

//Product struct preserves single DB row of product
type Product struct{
    ID int         `json:"id"`
    Model string   `json:"model"`
    Company string `json:"company"`
    Price int      `json:"price"`
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
    rows, err := db.Query("select * from products")
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    products := Products{}

    for rows.Next(){
        p:= Product {}
        err:= rows.Scan(&p.ID, &p.Model, &p.Company, &p.Price)
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
    req, err := db.Prepare("insert into products(model,company,price) values($1,$2,$3)")
    defer req.Close()
    if (err != nil) {
        return err
    }
    res, err := req.Exec(p.Model, p.Company, p.Price)
    if (err != nil) {
        return err
    }
    _, err = res.RowsAffected()
    if (err != nil) {
        return err
    }
    return nil
}
//DeleteProduct deletes product with selected ID from DB
func (db *DB) DeleteProduct(id int) error {
    req, err := db.Prepare("delete from products where id = ?")
    defer req.Close()
    if (err != nil) {
        return err
    }
    res, err := req.Exec(id)
    if (err != nil) {
        return err
    }
    _, err = res.RowsAffected()
    if (err != nil) {
        return err
    }
    return nil
}
