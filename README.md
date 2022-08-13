# go-sql-migrator

*Simple SQL migration utility for go, currently supports only mysql.*
*I decided to build this package to reduce the amount of code repetitions for various personal projects.*

### Author

* Lauri Orgla

### Supported features

* Inline sql declaration

```go
migrator.NewMigration(1, "Example inline declaration", func(db migrator.SQLClient) error {
			if _, err := db.Exec("Sql query here"); err != nil {
				return err
			}
			return nil
		})
```

* Embedded sql declaration via separate .sql files for each migration

```go
//go:embed relative_migrations_directory_path_here
var migrationsFS embed.FS
...
	migrations, err := migrator.CreateFromFS(migrationsFS)
	//db.DB is sql.DB compatible interface. sqlx, etc work as well. 
	m := migrator.NewMigrator(db.DB, log.Println)
	//Applies migrations, return error if not successful
	m.Apply(migrations)
```

### API
__Migrator__
* `migrator.NewMigrator(DatabaseClient, OptionalLoggingFunction) Migrator`
* `migrator.Apply([]Migrations) error`

__Migration__
* `migrator.NewMigration(id ,name, handler func(SQLClient) error) Migration`
* `migrator.CreateFromFS(embeddedFS) Migration`

### Examples

* [Embedded scripts](examples/embedded_scripts.go)
* [Explicit inline](examples/explicit_inline.go)