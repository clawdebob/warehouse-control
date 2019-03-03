package main

import (
    "./models"
    "fmt"
    "log"
    "net/http"
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

    http.HandleFunc("/products", env.getAllProducts)
    log.Fatal(http.ListenAndServe(":8080", nil))
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
    for _, p := range products {
        fmt.Fprintf(w, "%d, %s, %s, %d$\n", p.ID, p.Model,p.Company, p.Price)
    }
}
