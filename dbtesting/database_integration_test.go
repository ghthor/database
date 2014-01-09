package dbtesting

import (
	"github.com/ghthor/database/config"
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"github.com/ziutek/mymysql/mysql"
	"io/ioutil"
	"log"
)

var cfg config.Config

func init() {
	var err error
	cfg, err = config.ReadFromFile("../config.json")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
}

func checkIfDatabaseExists(c mysql.Conn, db string) (bool, error) {
	row, _, err := c.QueryFirst("select schema_name from information_schema.schemata where schema_name = '%s'", db)
	if err != nil {
		return false, err
	}

	return len(row) != 0, nil
}

func DescribeDatabaseIntegration(c gospec.Context) {
	// Create a Connection and Connect
	conn := mysql.New("tcp", "", "127.0.0.1:3306", cfg.Username, cfg.Password)
	c.Assume(conn.Connect(), IsNil)

	defer func() {
		err := conn.Close()
		c.Assume(err, IsNil)
	}()

	c.Specify("a test database", func() {
		c.Specify("can be created and dropped", func() {
			basename := "test-database"

			db, err := NewTestDatabase(basename, conn)
			c.Assume(err, IsNil)

			err = db.Create()
			c.Expect(err, IsNil)

			c.Specify("and is in use", func() {
				row, _, err := conn.QueryFirst("select DATABASE()")
				c.Assume(err, IsNil)
				c.Expect(row.Str(0), Equals, db.name)
			})

			dbExists, err := checkIfDatabaseExists(conn, db.name)
			c.Assume(err, IsNil)
			c.Expect(dbExists, IsTrue)

			err = db.Drop()
			c.Expect(err, IsNil)

			dbExists, err = checkIfDatabaseExists(conn, db.name)
			c.Assume(err, IsNil)
			c.Expect(dbExists, IsFalse)
		})

		c.Specify("can have a schema", func() {
			db, err := NewTestDatabase("test-database", conn)
			c.Assume(err, IsNil)

			c.Assume(db.Create(), IsNil)
			defer func() {
				c.Assume(db.Drop(), IsNil)
			}()

			schemaBytes, err := ioutil.ReadFile("test_schema.sql")
			c.Assume(err, IsNil)

			c.Assume(db.SetSchema(string(schemaBytes)), IsNil)

			_, err = db.Prepare("insert into test (name) values (?)")
			c.Expect(err, IsNil)

			c.Specify("that can only be set once", func() {
				c.Assume(db.SetSchema(""), Not(IsNil))
			})
		})

		c.Specify("generates a unique database name everytime", func() {
			basename := "unique-name-test"
			db1, err := NewTestDatabase(basename, conn)
			c.Assume(err, IsNil)

			db2, err := NewTestDatabase(basename, conn)
			c.Assume(err, IsNil)

			c.Expect(db1.name, Not(Equals), db2.name)
		})

		c.Specify("fails to create the database if a database using the name already exists", func() {
			genSuffix := func() (string, error) { return "non-unique", nil }
			basename := "failure-to-create"

			db1, err := newTestDatabase(basename, conn, genSuffix)
			c.Assume(err, IsNil)

			db2, err := newTestDatabase(basename, conn, genSuffix)
			c.Assume(err, IsNil)

			c.Assume(db1.Create(), IsNil)
			defer func() {
				c.Assume(db1.Drop(), IsNil)
			}()

			c.Expect(db2.Create(), Not(IsNil))
		})
	})
}
