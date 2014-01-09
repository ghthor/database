package database

import (
	"errors"
	"github.com/ghthor/database/action"
	"reflect"
)

var executorRegistry *ExecutorRegistry

func init() {
	executorRegistry = NewExecutorRegistry()
}

type NewExecutor func(DatabaseConn) (Executor, error)

type ExecutorRegistry struct {
	executors map[string]NewExecutor
}

func NewExecutorRegistry() *ExecutorRegistry {
	return &ExecutorRegistry{make(map[string]NewExecutor)}
}

func (r *ExecutorRegistry) Register(a action.A, e NewExecutor) error {
	typename := reflect.TypeOf(a).String()

	if _, exists := r.executors[typename]; exists {
		return errors.New("action binding already exists")
	} else {
		r.executors[typename] = e
	}

	return nil
}

func RegisterAction(a action.A, e NewExecutor) error {
	return executorRegistry.Register(a, e)
}

func (r *ExecutorRegistry) Lookup(a action.A) NewExecutor {
	typename := reflect.TypeOf(a).String()
	return r.executors[typename]
}

func (r *ExecutorRegistry) RegisteredActions() []action.A {
	return nil
}
