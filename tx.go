package database

import (
	"fmt"
	"github.com/ghthor/database/datatype"
	"github.com/ziutek/mymysql/mysql"
	"os"
	"path"
)

type RollbackError struct {
	err         error
	triggeredBy error
}

func (e RollbackError) Error() string {
	return fmt.Sprintf("%v after %v", e.err, e.triggeredBy)
}

type Transaction interface {
	Commit() error
	Rollback() error
	Run(mysql.Stmt, ...interface{}) (mysql.Result, error)
	//SaveFiles([]datatype.FormFile) ([]string, error)
	SaveFile(datatype.FormFile) (string, error)
}

type mysqlTransaction interface {
	Commit() error
	Rollback() error
	Do(mysql.Stmt) mysql.Stmt
}

type transaction struct {
	tx mysqlTransaction

	filepath   string
	savedFiles []string
}

func newTransaction(tx mysqlTransaction, filepath string) *transaction {
	return &transaction{tx, filepath, make([]string, 0, 1)}
}

func (t *transaction) Commit() error { return t.tx.Commit() }
func (t *transaction) Rollback() error {
	for _, filename := range t.savedFiles {
		err := os.Remove(path.Join(t.filepath, filename))
		if err != nil {
			return err
		}
	}
	return t.tx.Rollback()
}

func (t *transaction) Run(s mysql.Stmt, params ...interface{}) (mysql.Result, error) {
	// TODO: Specify params... with an Integration Test
	res, err := t.tx.Do(s).Run(params...)
	if err != nil {
		rollbackErr := t.Rollback()
		if rollbackErr != nil {
			return nil, RollbackError{rollbackErr, err}
		} else {
			return nil, err
		}
	}
	return res, nil
}

func (t *transaction) SaveFile(formFile datatype.FormFile) (string, error) {
	return t.saveFile(formFile, saveFile)
}

func (t *transaction) saveFile(formFile datatype.FormFile, savefn func(datatype.FormFile, string) (string, error)) (string, error) {
	filename, err := savefn(formFile, t.filepath)
	if err != nil {
		rollbackErr := t.Rollback()
		if rollbackErr != nil {
			return "", RollbackError{rollbackErr, err}
		} else {
			return "", err
		}
	}

	t.savedFiles = append(t.savedFiles, filename)

	return filename, nil
}
