package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/theorx/go-sql-migrator/pkg/migrator"
	"log"
)

func main() {
	log.Println("Running e2e test")

	//create sqlx client

	db, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_db")

	if err != nil {
		log.Println(err)
		return
	}

	if err := db.Ping(); err != nil {
		log.Println(err)
	}

	log.Println(db)

	m := migrator.NewMigrator(db.DB, log.Println)

	log.Println(m.Apply(nil))

}
