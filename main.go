package main

import (
    "./models"
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
)
//Env decribes struct that handles all requests
type Env struct {
    db models.Datastore
}

func main() {
    db, err := models.NewDB("test.db")
    if err != nil {
        log.Panic(err)
    }
    env := &Env{db}

    http.HandleFunc("/add/product", env.addProduct)
    http.HandleFunc("/products", env.getAllProducts)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func (env *Env) addProduct(w http.ResponseWriter, r *http.Request) {
    var p models.Product

    if r.Method != "POST" {
        http.Error(w, http.StatusText(405), 405)
        return
    }
    r.ParseForm()
    req, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal([]byte(req), &p)
    if (err != nil) {
        http.Error(w, http.StatusText(500), 500)
        return
    }
    err = env.db.InsertProduct(&p)
    if (err != nil) {
        http.Error(w, http.StatusText(500), 500)
        log.Fatal(err)
        return
    }
}

func (env *Env) getAllProducts(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    products, err := env.db.AllProducts()
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
    resp, _ := json.Marshal(products)
    fmt.Fprintln(w, string(resp))
}
