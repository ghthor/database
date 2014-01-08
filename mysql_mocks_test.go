package database

import (
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"github.com/ziutek/mymysql/mysql"
)

type MockMysqlConn struct {
	PrepareWasCalled bool
	PrepareFunc      func(string) (mysql.Stmt, error)
}

func (c *MockMysqlConn) Connect() error { return nil }
func (c *MockMysqlConn) Prepare(sql string) (mysql.Stmt, error) {
	c.PrepareWasCalled = true
	if c.PrepareFunc != nil {
		return c.PrepareFunc(sql)
	}
	return nil, nil
}
func (c *MockMysqlConn) Begin() (mysql.Transaction, error) { return nil, nil }

type MockStmt struct {
	RunWasCalled bool
	RunFunc      func(...interface{}) (mysql.Result, error)
}

func (s *MockStmt) Bind(params ...interface{}) {}

func (s *MockStmt) ResetParams() {}

func (s *MockStmt) Run(params ...interface{}) (mysql.Result, error) {
	s.RunWasCalled = true
	if s.RunFunc != nil {
		return s.RunFunc(params...)
	}
	return nil, nil
}
func (s *MockStmt) Delete() error { return nil }
func (s *MockStmt) Reset() error  { return nil }

func (s *MockStmt) SendLongData(pnum int, data interface{}, pkt_size int) error {
	return nil
}

func (s *MockStmt) Fields() []*mysql.Field { return nil }
func (s *MockStmt) NumField() int          { return 0 }
func (s *MockStmt) NumParam() int          { return 0 }
func (s *MockStmt) WarnCount() int         { return 0 }

func (s *MockStmt) Exec(params ...interface{}) ([]mysql.Row, mysql.Result, error) {
	return nil, nil, nil
}
func (s *MockStmt) ExecFirst(params ...interface{}) (mysql.Row, mysql.Result, error) {
	return nil, nil, nil
}
func (s *MockStmt) ExecLast(params ...interface{}) (mysql.Row, mysql.Result, error) {
	return nil, nil, nil
}

type MockResult struct {
	MessageStr string
}

func (r *MockResult) StatusOnly() bool           { return false }
func (r *MockResult) ScanRow(mysql.Row) error    { return nil }
func (r *MockResult) GetRow() (mysql.Row, error) { return nil, nil }

func (r *MockResult) MoreResults() bool                 { return false }
func (r *MockResult) NextResult() (mysql.Result, error) { return nil, nil }

func (r *MockResult) Fields() []*mysql.Field { return nil }
func (r *MockResult) Map(string) int         { return 0 }
func (r *MockResult) Message() string        { return r.MessageStr }
func (r *MockResult) AffectedRows() uint64   { return 0 }
func (r *MockResult) InsertId() uint64       { return 0 }
func (r *MockResult) WarnCount() int         { return 0 }

func (r *MockResult) MakeRow() mysql.Row              { return nil }
func (r *MockResult) GetRows() ([]mysql.Row, error)   { return nil, nil }
func (r *MockResult) End() error                      { return nil }
func (r *MockResult) GetFirstRow() (mysql.Row, error) { return nil, nil }
func (r *MockResult) GetLastRow() (mysql.Row, error)  { return nil, nil }

func DescribeMockMysqlConn(c gospec.Context) {
	// Compile time Verify interface implementation
	var _ MymysqlConn = &MockMysqlConn{}

	c.Specify("a mock sql conn", func() {
		conn := &MockMysqlConn{}
		c.Specify("can spy on Prepare method", func() {
			conn.Prepare("")
			c.Expect(conn.PrepareWasCalled, IsTrue)
		})

		c.Specify("can fake the Prepare implementation", func() {
			var argument string
			conn.PrepareFunc = func(sql string) (mysql.Stmt, error) {
				argument = sql
				return nil, nil
			}

			conn.Prepare("an sql statement")
			c.Expect(argument, Equals, "an sql statement")
		})
	})
}

func DescribeMockStmt(c gospec.Context) {
	// Static check for interface implementation
	var _ mysql.Stmt = &MockStmt{}

	c.Specify("a mock statement", func() {
		stmt := &MockStmt{}
		c.Specify("has a Run method", func() {
			c.Specify("that can be spied on", func() {
				res, err := stmt.Run()
				c.Expect(res, IsNil)
				c.Expect(err, IsNil)
				c.Expect(stmt.RunWasCalled, IsTrue)
			})

			c.Specify("that can be faked", func() {
				var arguments []interface{}

				stmt.RunFunc = func(params ...interface{}) (mysql.Result, error) {
					arguments = params
					return &MockResult{}, nil
				}

				params := []interface{}{1, "str"}
				res, err := stmt.Run(params...)

				_, isAMockResult := res.(*MockResult)
				c.Expect(isAMockResult, IsTrue)
				c.Expect(err, IsNil)
				c.Expect(stmt.RunWasCalled, IsTrue)
				c.Expect(len(arguments), Equals, len(params))
				for i, arg := range arguments {
					c.Expect(arg, Equals, params[i])
				}
			})
		})
	})
}

func DescribeMockResult() {
	// Static check for interface implementation
	var _ mysql.Result = &MockResult{}
}
