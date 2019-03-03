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

//AllProducts handles SQL request to get all products from DB
func (db *DB) AllProducts() ([]Product, error) {
    rows, err := db.Query("select * from products")
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    products := []Product{}

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
