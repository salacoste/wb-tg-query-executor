package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"gitlab.com/wb-dynamics/wb-tg-query-executor/internal/config"
	"gitlab.com/wb-dynamics/wb-tg-query-executor/internal/handler"
)

const (
    kCheckIntervalSec = 5
)

func main() {
    db := start()

    for {
        handler.DoWork(db)
        time.Sleep(time.Duration(kCheckIntervalSec) * time.Second)
    }
}

func start() *sql.DB {
    log.Printf("starting... ")
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not read config: %s", err)
	}

    log.Printf("connecting to database %s:%d user: %s dbname: %s", 
        conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Name)
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Name))
    if err != nil {
        log.Fatalf("could not connect to database: %s", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatalf("ping db failed: %s", err)
    }
    log.Printf("connected to database")
    return db
}

