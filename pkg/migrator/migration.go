package migrator

import (
	"embed"
	"golang.org/x/exp/maps"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

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

/*
CreateFromFS Reads all files that match pattern: "<integer>.sql" and orders them in ascending order
*/
func CreateFromFS(embeddedFS embed.FS) ([]Migration, error) {

	migrations := make(map[int]string)
	err := fs.WalkDir(embeddedFS, ".", func(path string, d fs.DirEntry, err error) error {
		//Ignore files that are not of .sql type
		if filepath.Ext(path) != ".sql" {
			return nil
		}

		nameWithoutExtension := strings.Split(filepath.Base(path), ".")[0]
		if num, err := strconv.Atoi(nameWithoutExtension); err == nil {
			migrations[num] = path
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	result := make([]Migration, 0)

	orderedKeys := maps.Keys(migrations)
	sort.Ints(orderedKeys)

	for _, index := range orderedKeys {
		fileName := migrations[index]
		result = append(result, NewMigration(int64(index), migrations[index], func(db SQLClient) error {

			f, err := embeddedFS.Open(fileName)
			if err != nil {
				return err
			}

			sqlBytes, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			if _, err := db.Exec(string(sqlBytes)); err != nil {
				return err
			}

			return nil
		}))
	}

	return result, nil
}
