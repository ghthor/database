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

	db, err := newDatabase(conn, filepath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// A Subset of the mysql.Conn interface to specify exactly what functionality we use
type mymysqlConn interface {
	Connect() error

	Prepare(string) (mysql.Stmt, error)
	Begin() (mysql.Transaction, error)
}

type database struct {
	mymysqlConn

	filepath   string
	fileServer http.Handler
}

func newDatabase(conn mymysqlConn, filepath string) (*database, error) {
	db := &database{
		mymysqlConn: conn,

		filepath:   filepath,
		fileServer: http.FileServer(http.Dir(filepath)),
	}

	return db, db.PrepareActions()
}

func (c *database) PrepareActions() (err error) {
	return
}

func (c *database) Execute(A action.A) (interface{}, error) {
	err := A.IsValid()
	if err != nil {
		return nil, err
	}

	fv := reflect.ValueOf(c).Elem().FieldByName(reflect.ValueOf(A).Type().Name())
	return fv.Interface().(Executor).ExecuteWith(A)
}
