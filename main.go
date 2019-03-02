package main

import (
    "./models"
    "context"
    "database/sql"
    "fmt"
    "log"
    "net/http"
)

type ContextInjector struct {
    ctx context.Context
    h   http.Handler
}

func (ci *ContextInjector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ci.h.ServeHTTP(w, r.WithContext(ci.ctx))
}

func main() {
    db, err := models.NewDB("test.db")
    if err != nil {
        log.Panic(err)
    }
    ctx := context.WithValue(context.Background(), "db", db)

    http.Handle("/products", &ContextInjector{ctx, http.HandlerFunc(booksIndex)})
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func booksIndex(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    db, ok := r.Context().Value("db").(*sql.DB)
    if !ok {
        http.Error(w, "could not get database connection pool from context", 500)
        return
    }

    products, err := models.AllProducts(db)
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
    for _, p := range products {
        fmt.Fprintf(w, "%d, %s, %s, %d$\n", p.ID, p.Model,p.Company, p.Price)
    }
}
