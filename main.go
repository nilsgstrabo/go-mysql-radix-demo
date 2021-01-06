package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/vrischmann/envconfig"
)

// Env from environment variables or arguments
type Env struct {
	Port   string `envconfig:"default=3003"`
	DbConn string `envconfig:""`
}

// Row defines a single mysql row
type Row struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// DataHandler implements http.Handler
type DataHandler struct {
	db *sql.DB
}

func main() {
	var env Env

	log.Info("starting")

	err := envconfig.Init(&env)
	if err != nil {
		log.Panic(err)
		panic(err)
	}

	db, err := sql.Open("mysql", env.DbConn)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	h := DataHandler{db: db}
	err = http.ListenAndServe(fmt.Sprintf(":%s", env.Port), h)
	if err != nil {
		log.Error(err)
	}
}

func (h DataHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	conn, err := h.db.Conn(ctx)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, "select id, name from demo.product")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	defer rows.Close()

	var (
		data []Row
		id   int32
		name string
	)
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
			return
		}
		data = append(data, Row{ID: id, Name: name})
	}

	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
