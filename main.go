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
    db, err := models.NewDB("test.db")
    if err != nil {
        log.Panic(err)
    }
    env := &Env{db}

    http.HandleFunc("/add/product", makeAddHandler(env.db.InsertProduct))
    http.HandleFunc("/delete/product", makeDeleteHandler(env.db.DeleteProduct))
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
        entities, err := fn()
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            log.Panic(err)
            return
        }
        resp, err := entities.ToJSON()
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            log.Panic(err)
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
                log.Panic(err)
                return
            }

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
            log.Panic(err)
            return
        }

    }
}
