package main

import (
    "./models"
    "log"
    "fmt"
    "net/http"
    "io/ioutil"
)
//Env decribes struct that contains handlers for all requests
type Env struct {
    db models.Datastore
}

func main() {
    const port = ":8080"
    db, err := models.NewDB("warehouse.db")
    if err != nil {
        log.Panic(err)
    }
    env := &Env{db}

    http.HandleFunc("/add/product", makeAddHandler(env.db.InsertProduct))
    http.HandleFunc("/add/person", makeAddHandler(env.db.InsertPerson))
    http.HandleFunc("/delete/product", makeDeleteHandler(env.db.DeleteProduct))
    http.HandleFunc("/delete/person", makeDeleteHandler(env.db.DeletePerson))
    http.HandleFunc("/products", makeGetAllHandler(env.db.AllProducts))
    http.HandleFunc("/persons", makeGetAllHandler(env.db.AllPersons))
    http.HandleFunc("/filter/product", makeFilterHandler(env.db.FilterProduct))
    log.Print("server has started on http://127.0.0.1" + port)
    log.Fatal(http.ListenAndServe(port, nil))
}

func makeGetAllHandler(fn func() (models.Serializable, error)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" {
            http.Error(w, http.StatusText(405), 405)
            return
        }
        entities, err := fn()
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            log.Print(err)
            return
        }
        resp, err := entities.ToJSON()
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            log.Print(err)
            return
        }
        fmt.Fprintln(w, resp)
    }
}

func makeAddHandler(fn func([]byte) error) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            if r.Method != "POST" {
                http.Error(w, http.StatusText(405), 405)
                return
            }
            r.ParseForm()
            req, _ := ioutil.ReadAll(r.Body)
            err := fn(req)
            if (err != nil) {
                http.Error(w, http.StatusText(500), 500)
                log.Print(err)
                return
            }

        }
}

func makeFilterHandler(fn func([]byte, string) (models.Serializable, error)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" {
            http.Error(w, http.StatusText(405), 405)
            return
        }
        r.ParseForm()
        req, _ := ioutil.ReadAll(r.Body)
        sort := r.Header.Get("sort")
        entities, err := fn(req, sort)
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            log.Print(err)
            return
        }
        resp, err := entities.ToJSON()
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            log.Print(err)
            return
        }
        fmt.Fprintln(w, resp)
    }
}

func makeDeleteHandler(fn func([]byte) error) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "DELETE" {
            http.Error(w, http.StatusText(405), 405)
            return
        }
        r.ParseForm()
        req, _ := ioutil.ReadAll(r.Body)
        err := fn(req)
        if (err != nil) {
            http.Error(w, http.StatusText(500), 500)
            log.Print(err)
            return
        }

    }
}
