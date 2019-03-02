package models

import (
    "fmt"
    "database/sql"
)

//product struct preserves single DB row of product
type product struct{
    ID int
    Model string
    Company string
    Price int
}

//AllProducts handles SQL request to get all products from DB
func AllProducts(db *sql.DB) ([]product, error) {
    rows, err := db.Query("select * from products")
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    products := []product{}

    for rows.Next(){
        p:= product {}
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
