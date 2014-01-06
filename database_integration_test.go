package database

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/ghthor/database/config"
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"github.com/ziutek/mymysql/mysql"
	"log"
)

var cfg config.Config

func init() {
	var err error
	cfg, err = config.ReadFromFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
}

func genSuffix() (string, error) {
	suffix := make([]byte, 16)
	n, err := rand.Read(suffix)
	if n != len(suffix) || err != nil {
		return "", err
	}

	return hex.EncodeToString(suffix), nil
}

func checkIfDatabaseExists(c mysql.Conn, db string) (bool, error) {
	row, _, err := c.QueryFirst("select schema_name from information_schema.schemata where schema_name = '%s'", db)
	if err != nil {
		return false, err
	}

	return len(row) != 0, nil
}

type TestDatabase struct {
	mysql.Conn
	name string
}

func NewTestDatabase(basename string, c mysql.Conn, genSuffix func() (string, error)) (*TestDatabase, error) {
	suffix, err := genSuffix()
	if err != nil {
		return nil, err
	}

	return &TestDatabase{c, basename + "_" + suffix}, nil
}

func (t *TestDatabase) Create() error {
	_, _, err := t.Conn.Query("CREATE DATABASE `%s` DEFAULT COLLATE = 'utf8_general_ci'", t.name)
	if err != nil {
		return err
	}
	return t.Use(t.name)
}

func (t *TestDatabase) Drop() error {
	_, _, err := t.Conn.Query("drop database `%s`", t.name)
	return err
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

			db, err := NewTestDatabase(basename, conn, genSuffix)
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

		c.Specify("generates a unique database name everytime", func() {
			basename := "unique-name-test"
			db1, err := NewTestDatabase(basename, conn, genSuffix)
			c.Assume(err, IsNil)

			db2, err := NewTestDatabase(basename, conn, genSuffix)
			c.Assume(err, IsNil)

			c.Expect(db1.name, Not(Equals), db2.name)
		})

		c.Specify("fails to create the database if a database using the name already exists", func() {
			genSuffix := func() (string, error) { return "non-unique", nil }
			basename := "failure-to-create"

			db1, err := NewTestDatabase(basename, conn, genSuffix)
			c.Assume(err, IsNil)

			db2, err := NewTestDatabase(basename, conn, genSuffix)
			c.Assume(err, IsNil)

			c.Assume(db1.Create(), IsNil)
			defer func() {
				c.Assume(db1.Drop(), IsNil)
			}()

			c.Expect(db2.Create(), Not(IsNil))
		})
	})

	c.Specify("an update statement", func() {
		db, err := NewTestDatabase("update-statement", conn, genSuffix)
		c.Assume(err, IsNil)

		err = db.Create()
		c.Assume(err, IsNil)

		defer func() {
			err := db.Drop()
			c.Assume(err, IsNil)
		}()

		res, err := db.Start(`
create table updateResultTest (
	id int AUTO_INCREMENT,
	txt text,
	PRIMARY KEY (id),
	UNIQUE KEY id (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

insert into updateResultTest (txt) values ('test');
`)
		c.Assume(err, IsNil)

		res, err = res.NextResult()
		c.Assume(err, IsNil)

		c.Specify("identifies the number of matching rows", func() {
			updateSql := "update updateResultTest set txt = 'updated' where id = %d limit 1"

			c.Specify("as 1", func() {
				res, err := db.Start(updateSql, 1)
				c.Assume(err, IsNil)

				updateResult := &UpdateResult{res}
				c.Expect(updateResult.MatchedRows(), Equals, uint64(1))
			})

			c.Specify("as none", func() {
				res, err := db.Start(updateSql, 2)
				c.Assume(err, IsNil)

				updateResult := &UpdateResult{res}
				c.Expect(updateResult.MatchedRows(), Equals, uint64(0))
			})
		})
	})
}
