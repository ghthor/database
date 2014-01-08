package database

import (
	"github.com/ghthor/database/config"
	"github.com/ghthor/database/dbtesting"
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

func DescribeDatabaseIntegration(c gospec.Context) {
	// Create a Connection and Connect
	conn := mysql.New("tcp", "", "127.0.0.1:3306", cfg.Username, cfg.Password)
	c.Assume(conn.Connect(), IsNil)

	defer func() {
		err := conn.Close()
		c.Assume(err, IsNil)
	}()

	c.Specify("an update statement", func() {
		db, err := dbtesting.NewTestDatabase("update-statement", conn)
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
