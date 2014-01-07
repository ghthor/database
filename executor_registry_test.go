package database

import (
	"errors"
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"github.com/ghthor/database/action"
	"reflect"
)

type (
	MockAction1 string
	MockAction2 string
	MockAction3 string
)

func (a MockAction1) IsValid() error { return errors.New(string(a)) }
func (a MockAction2) IsValid() error { return errors.New(string(a)) }
func (a MockAction3) IsValid() error { return errors.New(string(a)) }

func NewMockAction1Ex(DatabaseConn) (Executor, error) {
	return nil, nil
}

func NewMockAction2Ex(DatabaseConn) (Executor, error) {
	return nil, nil
}

func (e NewExecutor) Equals(other interface{}) bool {
	if fn, ok := other.(NewExecutor); ok {
		return reflect.ValueOf(e).Pointer() == reflect.ValueOf(fn).Pointer()
	}
	return false
}

func DescribeExecutorRegistry(c gospec.Context) {
	c.Specify("An executor registry", func() {
		r := NewExecutorRegistry()

		binds := []struct {
			action      action.A
			newExecutor NewExecutor
		}{
			{MockAction1(""), NewMockAction1Ex},
			{MockAction2(""), NewMockAction2Ex},
		}

		for _, bind := range binds {
			err := r.Register(bind.action, bind.newExecutor)
			c.Assume(err, IsNil)
		}

		c.Specify("can bind an action type to an executor constructor", func() {
			ex1 := r.Lookup(MockAction1(""))
			ex2 := r.Lookup(MockAction2(""))
			c.Expect(ex1, Equals, NewExecutor(NewMockAction1Ex))
			c.Expect(ex2, Equals, NewExecutor(NewMockAction2Ex))

			c.Specify("or will error if the action type has been bound already", func() {
				err := r.Register(MockAction1(""), NewMockAction2Ex)
				c.Assume(err, Not(IsNil))
				c.Expect(err.Error(), Equals, "action binding already exists")

				c.Specify("and will not modify the existing binding", func() {
					ex := r.Lookup(MockAction1(""))
					c.Expect(ex, Equals, NewExecutor(NewMockAction1Ex))
				})
			})
		})
	})
}
