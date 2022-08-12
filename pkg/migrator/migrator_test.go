package migrator

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"log"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewMigratorCreatesInstance(t *testing.T) {
	logCalled := false
	m := NewMigrator(nil, func(...any) {
		logCalled = true
	})

	m.log("test")
	if !logCalled {
		t.Errorf("failed to call correct log handler")
	}
}

func TestMigratorApplyAttemptsToMigrate(t *testing.T) {
	defer func() {
		//Clear hooks
		setupMigrationTableHook = nil
	}()

	testError := errors.New("result from setupMigrationTable")

	setupMigrationTableHook = func() error {
		return testError
	}

	m := NewMigrator(nil, nil)

	if m.Apply(nil) != testError {
		t.Errorf("Expected results from setupMigrationTableHook")
	}
}

func TestMigratorApplyAttemptsToFindMigrationID(t *testing.T) {
	defer func() {
		//Clear hooks
		setupMigrationTableHook = nil
		determineMigrationsStartedHook = nil
	}()
	setupMigrationTableHook = func() error { return nil }
	testError := errors.New("result from determineMigrationsStarted")

	determineMigrationsStartedHook = func() (int64, error) {
		return 0, testError
	}

	m := NewMigrator(nil, nil)

	if m.Apply(nil) != testError {
		t.Errorf("Expected error from determineMigrationsStarted")
	}
}

func TestMigratorApplyDeterminesIDAndLogsIt(t *testing.T) {
	defer func() {
		//Clear hooks
		setupMigrationTableHook = nil
		determineMigrationsStartedHook = nil
	}()
	setupMigrationTableHook = func() error { return nil }

	testId := int64(228)

	determineMigrationsStartedHook = func() (int64, error) {
		return testId, nil
	}

	logs := make([]string, 0)

	m := NewMigrator(nil, func(v ...any) {
		logs = append(logs, fmt.Sprintln(v...))
	})

	_ = m.Apply(nil)

	if !strings.Contains(logs[0], strconv.Itoa(int(testId))) {
		t.Errorf("Expected migration id in the log message")
	}
}

func TestMigratorApplyAppliesToMigrations(t *testing.T) {
	defer func() {
		//Clear hooks
		setupMigrationTableHook = nil
		determineMigrationsStartedHook = nil
		updateMigrationsAppliedHook = nil
	}()
	setupMigrationTableHook = func() error { return nil }
	updateMigrationsAppliedHook = func(Migration) error { return nil }

	lastMigrationID := int64(3)
	determineMigrationsStartedHook = func() (int64, error) {
		return lastMigrationID, nil
	}

	logs := make([]string, 0)

	m := NewMigrator(nil, func(v ...any) {
		logs = append(logs, fmt.Sprintln(v...))
	})

	applyCalls := 0
	dummyHandler := func(db SQLClient) error {
		applyCalls++
		return nil
	}

	migrations := []Migration{
		NewMigration(1, "first", dummyHandler),
		NewMigration(2, "second", dummyHandler),
		NewMigration(3, "third", dummyHandler),
		NewMigration(4, "fourth", dummyHandler),
		NewMigration(5, "fifth", dummyHandler),
	}

	//Expecting skip messages for 1,2,3 and calls for 4 and 5
	if err := m.Apply(migrations); err != nil {
		t.Errorf("Expected Apply() to return nil, got: %v", err)
	}

	if applyCalls != 2 {
		t.Errorf("Expected apply handlers to be called 2 times, got: %d", applyCalls)
	}

	for _, msg := range []string{
		"[Migrator]> Skipped ID: 1\n",
		"[Migrator]> Skipped ID: 2\n",
		"[Migrator]> Skipped ID: 3\n",
		"[Migrator]> Migration with id: 4 successfully applied!\n",
		"[Migrator]> Migration with id: 5 successfully applied!\n",
	} {
		if !slices.Contains(logs, msg) {
			t.Errorf("Log not found for: %s", msg)
		}
	}
}

func TestMigratorApplyFailsWithErrorFromMigration(t *testing.T) {
	defer func() {
		//Clear hooks
		setupMigrationTableHook = nil
		determineMigrationsStartedHook = nil
		updateMigrationsAppliedHook = nil
	}()
	setupMigrationTableHook = func() error { return nil }
	updateMigrationsAppliedHook = func(Migration) error { return nil }
	determineMigrationsStartedHook = func() (int64, error) { return 0, nil }

	m := NewMigrator(nil, nil)

	migrationError := errors.New("migration has failed")
	migrations := []Migration{
		NewMigration(12, "fault migration", func(db SQLClient) error {
			return migrationError
		}),
	}

	if got := m.Apply(migrations); got != migrationError {
		t.Errorf("Expected Apply() to return %v, got: %v", migrationError, got)
	}
}

func TestMigratorMigrateFailsUpgradeMigrationApplied(t *testing.T) {
	defer func() {
		updateMigrationsAppliedHook = nil
	}()

	testError := errors.New("test error")
	updateMigrationsAppliedHook = func(migration Migration) error {
		return testError
	}

	m := NewMigrator(nil, nil)
	entry := NewMigration(12, "fault migration", func(db SQLClient) error {
		return nil
	})

	if got := m.migrate(entry); got != testError {
		t.Errorf("migrate expected to return %v, got %v", got, got)
	}
}

func TestUpdateMigrationsAppliedExecFails(t *testing.T) {
	queryStringSeen := ""
	idSeen := int64(0)
	nameSeen := ""
	createdAtSeen := int64(0)

	testError := errors.New("exec error")

	m := NewMigrator(&mockSQLClient{
		execImpl: func(s string, a ...any) (*sql.Result, error) {
			queryStringSeen = s
			idSeen = a[0].(int64)
			nameSeen = a[1].(string)
			createdAtSeen = a[2].(int64)

			return nil, testError
		},
		queryRowImpl: nil,
	}, nil)

	entry := NewMigration(37, "x migration", func(db SQLClient) error {
		return nil
	})

	if got := m.updateMigrationsApplied(entry); got != testError {
		t.Errorf("updateMigrationsApplied expected to return %v, got %v", got, got)
	}

	if queryStringSeen != "INSERT INTO migrations (migration_id, name, created_at) VALUES(?, ?, ?)" {
		t.Errorf("incorrect query seen in exec")
	}

	if idSeen != 37 {
		t.Errorf("incorrect id seen in exec, expected 37, got: %d", idSeen)
	}

	if nameSeen != "x migration" {
		t.Errorf("incorrect name seen in exec, got: %s", nameSeen)
	}

	if time.Now().Unix()-createdAtSeen != 0 {
		t.Errorf("incorrect created_at seen in exec, time since created: %d", time.Now().Unix()-createdAtSeen)
	}
}

func TestUpdateMigrationsAppliedReturnsNil(t *testing.T) {
	m := NewMigrator(&mockSQLClient{
		execImpl: func(s string, a ...any) (*sql.Result, error) {

			return nil, nil
		},
		queryRowImpl: nil,
	}, nil)

	entry := NewMigration(12, "xyz migration", func(db SQLClient) error {
		return nil
	})

	if got := m.updateMigrationsApplied(entry); got != nil {
		t.Errorf("updateMigrationsApplied retured an error: %v, expected nil", got)
	}
}

func TestSetupMigrationTableExecsCorrectQuery(t *testing.T) {

	testError := errors.New("exec error")

	m := NewMigrator(&mockSQLClient{
		execImpl: func(s string, a ...any) (*sql.Result, error) {
			if s != "CREATE TABLE IF NOT EXISTS migrations (id INT AUTO_INCREMENT PRIMARY KEY,"+
				" migration_id INT NOT NULL UNIQUE, name varchar(128) NOT NULL, created_at INT NOT NULL) ENGINE=INNODB" {
				t.Errorf("Create table query did not match the expected query. Got: %s", s)
			}

			return nil, testError
		},
		queryRowImpl: nil,
	}, nil)

	if got := m.setupMigrationTable(); got != testError {
		t.Errorf("setupMigrationTable retured incorrect value. expected: %v, got: %v", testError, got)
	}
}

func TestSetupMigrationTableReturnsNil(t *testing.T) {
	m := NewMigrator(&mockSQLClient{
		execImpl: func(s string, a ...any) (*sql.Result, error) {
			return nil, nil
		},
		queryRowImpl: nil,
	}, nil)

	if got := m.setupMigrationTable(); got != nil {
		t.Errorf("setupMigrationTable returned error %v, expected nil", got)
	}
}

func TestDetermineMigrationStartingPointOnQueryErrorReturnsError(t *testing.T) {

	testError := errors.New("query row error")

	m := NewMigrator(&mockSQLClient{
		execImpl: nil,
		queryRowImpl: func(s string, a ...any) *sql.Row {

			return createRowWithError(testError)
		},
	}, nil)

	if _, got := m.determineMigrationStartingPoint(); got != testError {
		t.Errorf("determineMigrationStartingPoint returned error %v, expected %v", got, testError)
	}
}

/*
Mock implementation for testing db client interactions
*/
type mockSQLClient struct {
	execImpl     func(string, ...any) (*sql.Result, error)
	queryRowImpl func(string, ...any) *sql.Row
}

func (m *mockSQLClient) QueryRow(query string, args ...any) *sql.Row {
	return m.queryRowImpl(query, args...)
}

func (m *mockSQLClient) Exec(query string, args ...any) (*sql.Result, error) {
	return m.execImpl(query, args...)
}

func createRowWithError(err error) *sql.Row {
	row := &sql.Row{}
	//*(unsafe.Pointer(reflect.Indirect(reflect.ValueOf(row)).FieldByName("err").UnsafeAddr())) = err

	pointerVal := reflect.ValueOf(row)
	val := reflect.Indirect(pointerVal)

	log.Println(val.NumField())
	log.Println(val.Field(0))
	val.Field(1).Set(reflect.ValueOf(err).Elem())


	return row
}
