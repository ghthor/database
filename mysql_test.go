package database

import (
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"github.com/ziutek/mymysql/mysql"
	"io/ioutil"
)

func DescribeMysqlDatabaseIntegration(c gospec.Context) {
	// Create a Connection and Connect
	conn := mysql.New("tcp", "", "127.0.0.1:3306", cfg.Username, cfg.Password)
	c.Assume(conn.Connect(), IsNil)

	defer func() {
		err := conn.Close()
		c.Assume(err, IsNil)
	}()

	c.Specify("a mysql database", func() {
		c.Specify("can be created and dropped", func() {
			basename := "test-database"

			db, err := NewUniqMysqlDatabase(basename, conn)
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

		schemaBytes, err := ioutil.ReadFile("dbtesting/test_schema.sql")
		c.Assume(err, IsNil)

		c.Specify("can have a schema", func() {
			db, err := NewUniqMysqlDatabase("test-database", conn)
			c.Assume(err, IsNil)

			c.Assume(db.Create(), IsNil)
			defer func() {
				c.Assume(db.Drop(), IsNil)
			}()

			c.Assume(db.SetSchema(string(schemaBytes)), IsNil)

			_, err = db.Prepare("insert into test (name) values (?)")
			c.Expect(err, IsNil)

			c.Specify("that can only be set once", func() {
				c.Assume(db.SetSchema(""), Not(IsNil))
			})
		})

		c.Specify("will use the existing database", func() {
			db, err := NewUniqMysqlDatabase("test-database", conn)
			c.Assume(err, IsNil)

			c.Assume(db.Create(), IsNil)
			defer func() {
				c.Assume(db.Drop(), Not(IsNil))
			}()

			c.Assume(db.SetSchema(string(schemaBytes)), IsNil)

			conn2 := mysql.New("tcp", "", "127.0.0.1:3306", cfg.Username, cfg.Password)
			c.Assume(conn2.Connect(), IsNil)

			defer func() {
				err := conn2.Close()
				c.Assume(err, IsNil)
			}()

			dbReused, err := NewMysqlDatabase(db.name, conn2)
			c.Assume(err, IsNil)

			defer func() {
				err := dbReused.Drop()
				c.Assume(err, IsNil)
			}()

			c.Expect(dbReused.Create(), Not(IsNil))
			row, _, err := conn2.QueryFirst("select DATABASE()")
			c.Assume(err, IsNil)
			c.Expect(row.Str(0), Equals, dbReused.name)

			c.Specify("and cannot set schema", func() {
				c.Expect(dbReused.SetSchema(string(schemaBytes)), Not(IsNil))
			})
		})

		c.Specify("can have a unique name", func() {
			basename := "unique-name-test"
			db1, err := NewUniqMysqlDatabase(basename, conn)
			c.Assume(err, IsNil)

			db2, err := NewUniqMysqlDatabase(basename, conn)
			c.Assume(err, IsNil)

			c.Expect(db1.name, Not(Equals), db2.name)
		})

		c.Specify("fails to create the database if a database using the name already exists", func() {
			genSuffix := func() (string, error) { return "non-unique", nil }
			basename := "failure-to-create"

			db1, err := newMysqlDatabase(basename, conn, genSuffix)
			c.Assume(err, IsNil)

			db2, err := newMysqlDatabase(basename, conn, genSuffix)
			c.Assume(err, IsNil)

			c.Assume(db1.Create(), IsNil)
			defer func() {
				c.Assume(db1.Drop(), IsNil)
			}()

			c.Expect(db2.Create(), Not(IsNil))
		})
	})

	c.Specify("an update statement", func() {
		db, err := NewUniqMysqlDatabase("update-statement", conn)
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
