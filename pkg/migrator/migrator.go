package migrator

import (
	"database/sql"
	"time"
)

// Handle schema migrations
type migrator struct {
	db  SQLClient
	log LoggerFunction
}

/*
NewMigrator Creates instance of migrator to apply migrations
If logging is not required, just pass nil as logFunction Otherwise a good substitute could be log.Println
*/
func NewMigrator(database SQLClient, logFunction LoggerFunction) *migrator {
	return &migrator{db: database, log: func(v ...any) {
		if logFunction != nil {
			logFunction(v...)
		}
	}}
}

/*
Apply migrations based on database state
All migrations are only applied once, consecutive runs should be safe as already applied migrations
are ignored.
*/
func (m *migrator) Apply(migrations []Migration) error {
	if err := m.setupMigrationTable(); err != nil {
		return err
	}

	migrationID, err := m.determineMigrationStartingPoint()
	if err != nil {
		return err
	}

	m.log("[Migrator]> Applying migrations.. starting from:", migrationID)
	for _, entry := range migrations {
		if entry.ID() <= migrationID {
			m.log("[Migrator]> Skipped ID:", entry.ID())
			continue
		}

		if err := m.migrate(entry); err != nil {
			return err
		}
	}

	m.log("[Migrator]> All finished successfully")
	return nil
}

func (m *migrator) migrate(entry Migration) error {
	m.log("[Migrator]> Updating database for:", entry.ID(), "-", entry.Name())
	if err := m.updateMigrationsApplied(entry); err != nil {
		m.log("[Migrator]> Database update has failed, aborting migration! Error:", err)
		return err
	}

	m.log("[Migrator]> Applying migration for:", entry.ID(), "-", entry.Name())
	if err := entry.Apply(m.db); err != nil {
		m.log("[Migrator]> Applying migration has failed, aborting migration! Error:", err, "Entry:", entry)
		return err
	}

	m.log("[Migrator]> Migration with id:", entry.ID(), "successfully applied!")
	return nil
}

func (m *migrator) determineMigrationStartingPoint() (int64, error) {
	if determineMigrationsStartedHook != nil {
		return determineMigrationsStartedHook()
	}
	m.log("[Migrator]> Determining the latest migration applied to the database..")

	row := m.db.QueryRow("SELECT MAX(migration_id) as 'last_migration' FROM migrations")
	if err := row.Err(); err != nil {
		return 0, err
	}

	selectField := &sql.NullInt64{}

	if err := row.Scan(&selectField); err != nil {
		return 0, err
	}

	if selectField == nil {
		m.log("[Migrator]> Migration id was not found, starting from 0")
		return 0, nil
	}

	return selectField.Int64, nil
}

/*
Creates table to a mysql database that is currently in use by the provided connection
*/
func (m *migrator) setupMigrationTable() error {
	if setupMigrationTableHook != nil {
		return setupMigrationTableHook()
	}

	m.log("[Migrator]> Setting up migrations table")

	if _, err := m.db.Exec(
		"CREATE TABLE IF NOT EXISTS migrations (id INT AUTO_INCREMENT PRIMARY KEY," +
			" migration_id INT NOT NULL UNIQUE, name varchar(128) NOT NULL, created_at INT NOT NULL) ENGINE=INNODB",
	); err != nil {
		m.log("[Migrator]> migration table initialization has failed, error:", err)
		return err
	}

	return nil
}

/*
Updates migrations table with the given entry information
*/
func (m *migrator) updateMigrationsApplied(entry Migration) error {
	if updateMigrationsAppliedHook != nil {
		return updateMigrationsAppliedHook(entry)
	}

	if _, err := m.db.Exec(
		"INSERT INTO migrations (migration_id, name, created_at) VALUES(?, ?, ?)",
		entry.ID(),
		entry.Name(),
		time.Now().Unix()); err != nil {
		return err
	}

	return nil
}

type LoggerFunction func(v ...any)

// Hooks for testing
var setupMigrationTableHook func() error
var updateMigrationsAppliedHook func(Migration) error
var determineMigrationsStartedHook func() (int64, error)

type SQLClient interface {
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}
