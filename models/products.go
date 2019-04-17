package models

import (
    "fmt"
    "encoding/json"
    "crypto/rand"
    "strings"
)

//Product represents single DB row of product
type Product struct{
    Serial string       `json:"serial"`
    Name string         `json:"name"`
    Company string      `json:"company"`
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

func (db *DB) InitRef() {
    db.constructFilterQuery("dsf", Product{"ffff", "dank", "memes", 1, 1, 1});
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
    finalQuery := "UPDATE Goods SET"
    query := make([]string, 0)
    err := json.Unmarshal(data, &p)
    if (err != nil) {
        return err
    }
    if p.Name != "" {
        query = append(query, fmt.Sprintf(" Name = '%s'", p.Name))
    }
    if p.Company != "" {
        query = append(query, fmt.Sprintf(" Manufacturer = '%s'", p.Company))
    }
    if p.Place != 0 {
        query = append(query, fmt.Sprintf(" Place = %d", p.Place))
    }
    if p.Row != 0 {
        query = append(query, fmt.Sprintf(" Row = %d", p.Row))
    }
    if p.Column != 0 {
        query = append(query, fmt.Sprintf(" Columm = %d", p.Column))
    }
    finalQuery += strings.Join(query, " ,")
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
    if (err != nil) {
        return err
    }
    id := p.Serial
    return db.execEntity("DELETE FROM Goods WHERE Serial = ?", id)
}
//FilterProduct filters Good's rows according to specified filters
func (db *DB) FilterProduct(data []byte, sort string) (Serializable, error) {
    var p Product
    sortBy := " ORDER BY Name"
    query := make([]string, 0)
    err := json.Unmarshal(data, &p)
    if (err != nil) {
        return nil, err
    }
    finalQuery, err := db.constructFilterQuery("Goods", p)
    if (err != nil) {
        return nil, err
    }
    if sort == "desc" {
        sortBy += " DESC"
    }
    finalQuery += strings.Join(query, " AND") + sortBy
    return db.getProductsQuery(finalQuery)
}
