package models

import(
  "fmt"
  "encoding/json"
)

type Order struct{
  Id uint
  Serial string
  Date string
  Type int
  ClinetId int
}

type Orders []Order

func (o Orders) ToJSON() (string, error) {
  json, err := json.MarshalIndent(o, "", "    ")
  return string(json), err
}

func (db *DB) getOrdersQuerry(q string) (Serializable, error) {
  fmt.Println(q);
  rows, err := db.Query(q)
  if err != nil {
    return nil, err
  }

  defer rows.Close()
  orders := Orders{}

  for rows.Next() {
    o := Order {}
    err := rows.Scan(
      &o.Id,
      &o.Serial,
      &o.Date,
      &o.Type,
      &o.ClinetId,
    )

    if err != nil {
      return nil, err
    }
    orders = append(orders, o)
  }
  return orders, nil
}

func (db *DB) AllOrders() (Serializable, error){
  return db.getOrdersQuerry("SELECT * FROM Orders")
}

func (db *DB) InsertOrder(parse []byte) error {
  var o Order
  err := json.Unmarshal(parse, &o)
  if(err != nil){
    return nil
  }

  return db.execEntity(
    "INSERT INTO Orders(Id, Serial, Date, Type, ClientId) VALUES($1,$2,$3,$4,$5)",
    o.Id,
    o.Serial,
    o.Date,
    o.Type,
    o.ClinetId,
  )
}
