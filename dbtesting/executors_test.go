package dbtesting

import (
	"github.com/ghthor/database/action"
	"github.com/ghthor/database/config"
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"io/ioutil"
	"log"
	"os"
)

var cfg config.Config

func init() {
	var err error
	cfg, err = config.ReadFromFile("../config.json")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
}

type MockAction struct{}

func (MockAction) IsValid() error { return nil }

type MockExecutor struct {
	c *ExecutorContext

	connWasOpened bool
	dbWasCreated  bool
	dirWasCreated bool
}

func (MockExecutor) ExecuteWith(a action.A) (interface{}, error) {
	return nil, nil
}

func (e *MockExecutor) Describe(c *ExecutorContext) {
	var err error
	e.c = c

	e.connWasOpened = c.Db.MysqlDatabase().IsConnected()

	e.dbWasCreated, err = c.Db.MysqlDatabase().Exists()
	c.Assume(err, IsNil)

	_, err = os.Open(c.Db.Filepath())

	e.dirWasCreated = (err == nil)
}

func DescribeSpecifyExecutor(c gospec.Context) {
	schemaBytes, err := ioutil.ReadFile("test_schema.sql")
	c.Assume(err, IsNil)

	c.Specify("an executor context", func() {
		executor := &MockExecutor{}

		var databaseWasDeleted bool

		DescribeExecutor(c, MockAction{}, executor, cfg, string(schemaBytes), func() {
			exists, err := executor.c.Db.MysqlDatabase().Exists()
			c.Assume(err, IsNil)

			databaseWasDeleted = !exists
		})
		c.Assume(executor.c, Not(IsNil))

		c.Specify("opens an mysql connection", func() {
			c.Expect(executor.connWasOpened, IsTrue)

			c.Specify("and closes it after upon completion", func() {
				c.Expect(executor.c.Db.MysqlDatabase().IsConnected(), IsFalse)
			})
		})

		c.Specify("creates a new database", func() {
			c.Expect(executor.dbWasCreated, IsTrue)

			c.Specify("and removes it upon completion", func() {
				c.Expect(databaseWasDeleted, IsTrue)
			})
		})

		c.Specify("creates a temporary directory", func() {
			c.Expect(executor.dirWasCreated, IsTrue)

			c.Specify("and removes it upon completion", func() {
				_, err = os.Open(executor.c.Db.Filepath())
				c.Expect(err, Not(IsNil))
				_, IsPathError := err.(*os.PathError)
				c.Expect(IsPathError, IsTrue)
			})
		})
	})
}
