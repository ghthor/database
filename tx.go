package database

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
)

type RollbackError struct {
	err         error
	triggeredBy error
}

func (e RollbackError) Error() string {
	return fmt.Sprintf("%v after %v", e.err, e.triggeredBy)
}

func TxStmtRun(tr mysql.Transaction, stmt mysql.Stmt, params ...interface{}) (mysql.Result, error) {
	res, err := tr.Do(stmt).Run(params...)
	if err != nil {
		rollbackErr := tr.Rollback()
		if rollbackErr != nil {
			return nil, RollbackError{rollbackErr, err}
		} else {
			return nil, err
		}
	}
	return res, nil
}
