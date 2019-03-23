package main

import (
    "./models"
    "log"
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
)
//Env decribes struct that handles all requests
type Env struct {
    db models.Datastore
}

func main() {
    const port = ":8080"
    db, err := models.NewDB("test.db")
    if err != nil {
        log.Panic(err)
    }
    env := &Env{db}

    http.HandleFunc("/add/product", env.addProduct)
    http.HandleFunc("/delete/product", env.deleteProduct)
    http.HandleFunc("/products", makeGetAllHandler(env.db.AllProducts))
    log.Print("server has started on http://127.0.0.1" + port)
    log.Fatal(http.ListenAndServe(port, nil))
}

func makeGetAllHandler(fn func() (models.Serializable, error)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" {
            http.Error(w, http.StatusText(405), 405)
            return
        }
        products, err := fn()
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            log.Panic(err)
            return
        }
        resp, err := products.ToJSON()
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            log.Panic(err)
            return
        }
        fmt.Fprintln(w, resp)
    }
}

func (env *Env) addProduct(w http.ResponseWriter, r *http.Request) {
    var p models.Product

    if r.Method != "POST" {
        http.Error(w, http.StatusText(405), 405)
        return
    }
    r.ParseForm()
    req, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal(req, &p)
    if (err != nil) {
        http.Error(w, http.StatusText(500), 500)
        return
    }
    err = env.db.InsertProduct(&p)
    if (err != nil) {
        http.Error(w, http.StatusText(500), 500)
        log.Panic(err)
        return
    }
}

func (env *Env) deleteProduct(w http.ResponseWriter, r *http.Request) {
    var p models.Product

    if r.Method != "DELETE" {
        http.Error(w, http.StatusText(405), 405)
        return
    }
    r.ParseForm()
    req, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal(req, &p)
    if (err != nil) {
        http.Error(w, http.StatusText(500), 500)
        log.Panic(err)
        return
    }
    err = env.db.DeleteProduct(p.ID)
    if (err != nil) {
        http.Error(w, http.StatusText(500), 500)
        log.Fatal(err)
        return
    }
}
