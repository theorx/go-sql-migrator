package migrator

import (
	"errors"
	"testing"
)

func TestNewMigrationCreatesInstance(t *testing.T) {

	handlerFn := func(db SQLClient) error {
		return errors.New("example-handler")
	}

	m := NewMigration(4, "test-4", handlerFn)

	if m.ID() != 4 {
		t.Errorf("Expected migration id to be 4, got: %d", m.ID())
	}

	if m.Name() != "test-4" {
		t.Errorf("Expected migration id to be 'test-4', got: %s", m.Name())
	}

	if m.Apply(nil).Error() != "example-handler" {
		t.Errorf("Apply did not call the correct function")
	}
}
