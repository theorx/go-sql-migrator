package main

import (
	"embed"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/theorx/go-sql-migrator/pkg/migrator"
	"log"
)

//go:embed embedded_migrations
var migrationsFS embed.FS

/*
Embedded scripts example
*/
func main() {
	db, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_db")

	if err != nil {
		log.Println(err)
		return
	}

	/*
		Using the embedded filesystem, the CreateFromFS will crawl and find all sql files which match
		"<integer>.sql" format, then orders them in ascending order and constructs the migrations list based on that.
	*/
	migrations, err := migrator.CreateFromFS(migrationsFS)
	m := migrator.NewMigrator(db.DB, log.Println)

	log.Println(m.Apply(
		migrations,
	))
}
