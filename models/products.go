package models

import (
    "fmt"
    "encoding/json"
    "crypto/rand"
    "strings"
    "strconv"
)

//Product represents single DB row of product
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
    return db.getProductsQuery("select * from Goods")
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
        "insert into Goods values($1,$2,$3,$4,$5,$6)",
        p.Serial,
        p.Name,
        p.Company,
        p.Row,
        p.Column,
        p.Place,
    )
}
//DeleteProduct deletes product with selected Serial from DB
func (db *DB) DeleteProduct(parse []byte) error {
    var p Product
    err := json.Unmarshal(parse, &p)
    if (err != nil) {
        return err
    }
    id := p.Serial
    return db.execEntity("delete from Goods where Serial = ?", id)
}
//FilterProduct filters Good's rows according to specified filters
func (db *DB) FilterProduct(data []byte, sort string) (Serializable, error) {
    var p Product
    finalQuery := "Select * from Goods where"
    sortBy := " order by Name"
    query := make([]string, 0)
    err := json.Unmarshal(data, &p)
    if (err != nil) {
        return nil, err
    }
    if p.Serial != "" {
        query = append(query, " Serial like '" + p.Serial + "%'")
    }
    if p.Name != "" {
        query = append(query, " Name like '" + p.Name + "%'")
    }
    if p.Company != "" {
        query = append(query, " Manufacturer like '" + p.Company + "%'")
    }
    if p.Place != 0 {
        query = append(query, " Place = " + strconv.Itoa(p.Place))
    }
    if p.Row != 0 {
        query = append(query, " Row = " + strconv.Itoa(p.Row))
    }
    if p.Column != 0 {
        query = append(query, " Columm = " + strconv.Itoa(p.Column))
    }
    if sort == "desc" {
        sortBy += " desc"
    }
    finalQuery += strings.Join(query, " and") + sortBy
    return db.getProductsQuery(finalQuery)
}
