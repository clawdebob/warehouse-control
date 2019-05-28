package models

import(
  "fmt"
  "encoding/json"
)

type Order struct{
  Id uint           `json:"id"`
  Serial string     `json:"serial"`
  Date string       `json:"date"`
  Type int          `json:"type"`
  ClinetId int      `json:"client_id"`
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
    "INSERT INTO Orders(Serial, Date, Type, ClientId) VALUES($1,$2,$3,$4)",
    o.Serial,
    o.Date,
    o.Type,
    o.ClinetId,
  )
}

func (db *DB)DeleteOrder(parse []byte) error{
  var o Order
  err := json.Unmarshal(parse, &o)
  if (err != nil) {
    return err
  }
  return db.execEntity("DELETE FROM Orders WHERE Id = ?", o.Id)
}

func (db* DB) EditOrder(data []byte) error{
  var o Order
  err := json.Unmarshal(data, &o)
  if (err != nil){
    return err
  }

  finalQuery, err := db.update("Orders", o)
  if (err != nil){
    return err
  }

  if o.Id != 0 {
    finalQuery += fmt.Sprintf(" WHERE Id = %d", o.Id)
  } else {
    return fmt.Errorf("id is invalid")
  }

  return db.execEntity(finalQuery)
}

func (db *DB) FilterOrder(data []byte, sort string) (Serializable, error){
	var o Order
	sortBy := " ORDER BY Date"
	err:= json.Unmarshal(data, &o)
	if (err != nil){
		return nil, err
	}

	finalQuery, err :=  db.filter("Orders", o)
	if (err != nil){
		return nil, err
	}

	if (sort = "desc"){
		sortBy+=" DESC"
	}

	finalQuery+=sortBy
	return db.GetOrdersQuery(finalQuery)
}
