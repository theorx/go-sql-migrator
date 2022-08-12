package migrator

type Migration interface {
	Name() string
	ID() int64
	Apply(SQLClient) error
}

type migrationImpl struct {
	id        int64
	name      string
	applyFunc func(db SQLClient) error
}

func (m *migrationImpl) Name() string {
	return m.name
}

func (m *migrationImpl) ID() int64 {
	return m.id
}

func (m *migrationImpl) Apply(client SQLClient) error {
	return m.applyFunc(client)
}

/*
NewMigration Creates new migration
*/
func NewMigration(id int64, name string, handler func(db SQLClient) error) Migration {
	return &migrationImpl{
		id:        id,
		name:      name,
		applyFunc: handler,
	}
}
