package dbtesting

import (
	"fmt"
	"github.com/ghthor/database"
	"github.com/ghthor/database/action"
	"github.com/ghthor/database/config"
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"github.com/ziutek/mymysql/mysql"
	"io/ioutil"
	"os"
	"reflect"
)

type ExecutorDescription interface {
	Describe(*ExecutorContext)
}

type ExecutorContext struct {
	gospec.Context

	Conn database.MymysqlConn
	Db   *database.Database

	Input action.A
	Impl  database.Executor

	Res interface{}
	Err error
}

func (c *ExecutorContext) Run() {
	c.Res, c.Err = c.Impl.ExecuteWith(c.Input)

}

func (c *ExecutorContext) SpecifyResult(expectedResult interface{}) {
	c.Specify(fmt.Sprintf("should return a [%s]", reflect.ValueOf(expectedResult).Type()), func() {
		c.Run()
		c.Expect(c.Res, Equals, expectedResult)
		c.Expect(c.Err, IsNil)
	})
}

func (c *ExecutorContext) SpecifySideEffects(description string, expectations func()) {
	c.Specify(description, func() {
		c.Run()
		expectations()
	})
}

func DescribeExecutor(c gospec.Context, input action.A, e ExecutorDescription, cfg config.Config, schema string, beforeClose func()) {
	c.Specify(fmt.Sprintf("the [%s] action should be executed by [%s] and", reflect.TypeOf(input), reflect.TypeOf(e)), func() {
		describeExecutor(c, input, e, cfg, schema, beforeClose)
	})
}

func describeExecutor(c gospec.Context, input action.A, e ExecutorDescription, cfg config.Config, schema string, beforeClose func()) {
	conn := mysql.New("tcp", "", "127.0.0.1:3306", cfg.Username, cfg.Password)
	c.Assume(conn.Connect(), IsNil)

	defer func() {
		if beforeClose != nil {
			beforeClose()
		}

		err := conn.Close()
		c.Assume(err, IsNil)
	}()

	testDb, err := database.NewUniqMysqlDatabase(cfg.DefaultDB, conn)
	c.Assume(err, IsNil)

	err = testDb.Create()
	c.Assume(err, IsNil)

	defer func() {
		err := testDb.Drop()
		c.Assume(err, IsNil)
	}()

	c.Assume(testDb.SetSchema(schema), IsNil)

	tmp, err := ioutil.TempDir("", cfg.DefaultDB)
	c.Assume(err, IsNil)

	defer func() {
		c.Assume(os.RemoveAll(tmp), IsNil)
	}()

	db, err := database.NewDatabase(testDb, tmp)
	c.Assume(err, IsNil)

	context := &ExecutorContext{
		Context: c,

		Db: db,

		Input: input,
		Impl:  e.(database.Executor),
	}

	e.Describe(context)
}
