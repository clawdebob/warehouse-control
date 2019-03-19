package models

import (
    "fmt"
)

//Product struct preserves single DB row of product
type Product struct{
    ID int
    Model string
    Company string
    Price int
}

//Products is a type that describes Products array
type Products []Product

//AllProducts handles SQL request to get all products from DB
func (db *DB) AllProducts() (Products, error) {
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

    for _, p := range products {
        fmt.Println(p.ID, p.Model, p.Company, p.Price)
    }
    return products, nil
}
//InsertProduct adds a new product entry in DB
func (db *DB) InsertProduct(p *Product) error {
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
