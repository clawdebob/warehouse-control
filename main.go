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

    http.HandleFunc("/add/product", makeTxHandler(env.db.InsertProduct))
    http.HandleFunc("/add/person", makeTxHandler(env.db.InsertPerson))
    http.HandleFunc("/delete/product", makeDeleteHandler(env.db.DeleteProduct))
    http.HandleFunc("/delete/person", makeDeleteHandler(env.db.DeletePerson))
    http.HandleFunc("/products", makeGetAllHandler(env.db.AllProducts))
    http.HandleFunc("/persons", makeGetAllHandler(env.db.AllPersons))
    http.HandleFunc("/filter/product", makeFilterHandler(env.db.FilterProduct))
    http.HandleFunc("/filter/person", makeFilterHandler(env.db.FilterPerson))
    http.HandleFunc("/edit/product", makeTxHandler(env.db.EditProduct))
    http.HandleFunc("/edit/person", makeTxHandler(env.db.EditPerson))
    log.Print("server has started on http://127.0.0.1" + port)
    log.Fatal(http.ListenAndServe(port, nil))
}

//Enable CORS policy in chrom just in case
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func makeDeleteHandler(fn func([]byte) error) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w)
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


func makeGetAllHandler(fn func() (models.Serializable, error)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w)
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

func makeTxHandler(fn func([]byte) error) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            keys := make([]string, 0)
            enableCors(&w)
            if r.Method != "POST" {
                http.Error(w, http.StatusText(405), 405)
                return
            }
            r.ParseForm()
            for key := range r.Form {
                keys = append(keys, key)
            }
            err := fn([]byte(keys[0]))
            if (err != nil) {
                http.Error(w, http.StatusText(500), 500)
                log.Print(err)
                return
            }

        }
}

func makeFilterHandler(fn func([]byte, string) (models.Serializable, error)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        keys := make([]string, 0)
        enableCors(&w)
        if r.Method != "GET" {
            http.Error(w, http.StatusText(405), 405)
            return
        }
        r.ParseForm()
        sort := r.Header.Get("sort")
        for key := range r.Form {
            keys = append(keys, key)
            fmt.Println(key)
        }
        entities, err := fn([]byte(keys[0]), sort)
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
