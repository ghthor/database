package database

import (
	"github.com/ghthor/database/action"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	"net/http"
	"reflect"
)

type Db interface {
	Execute(action.A) (interface{}, error)
}

type Executor interface {
	ExecuteWith(action.A) (interface{}, error)
}

func New(user, passwd, database, filepath string) (Db, error) {
	conn := mysql.New("tcp", "", "127.0.0.1:3306", user, passwd, database)
	err := conn.Connect()
	if err != nil {
		return nil, err
	}

	mysqlDb, err := NewMysqlDatabase(database, conn)
	if err != nil {
		return nil, err
	}

	db, err := NewDatabase(mysqlDb, filepath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// A Subset of the mysql.Conn interface to specify exactly what functionality we use
type MymysqlConn interface {
	Connect() error

	Prepare(string) (mysql.Stmt, error)
	Begin() (mysql.Transaction, error)
}

type DatabaseConn interface {
	MysqlConn() MymysqlConn
	Filepath() string
	Begin() (Transaction, error)
}

type Database struct {
	mysqlDb *MysqlDatabase

	filepath   string
	fileServer http.Handler
}

func NewDatabase(mysqlDb *MysqlDatabase, filepath string) (*Database, error) {
	db := &Database{
		mysqlDb: mysqlDb,

		filepath:   filepath,
		fileServer: http.FileServer(http.Dir(filepath)),
	}

	return db, db.PrepareActions()
}

func (c *Database) MysqlConn() MymysqlConn { return c.mysqlDb }
func (c *Database) Filepath() string       { return c.filepath }
func (c *Database) Begin() (Transaction, error) {
	tx, err := c.MysqlConn().Begin()
	if err != nil {
		return nil, err
	}
	return newTransaction(tx, c.filepath), nil
}

func (c *Database) MysqlDatabase() *MysqlDatabase { return c.mysqlDb }

func (c *Database) PrepareActions() (err error) {
	return
}

func (c *Database) Execute(A action.A) (interface{}, error) {
	err := A.IsValid()
	if err != nil {
		return nil, err
	}

	fv := reflect.ValueOf(c).Elem().FieldByName(reflect.ValueOf(A).Type().Name())
	return fv.Interface().(Executor).ExecuteWith(A)
}
