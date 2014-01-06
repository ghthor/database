package database

import (
	"errors"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
)

type UpdateResult struct {
	mysql.Result
}

var ErrUnexpectedMessageFormat = errors.New("unexpected message format")

func (r *UpdateResult) MatchedRows() uint64 {
	var matched uint64
	n, err := fmt.Sscanf(r.Message(), "(Rows matched: %d Changed:", &matched)
	if n != 1 || err != nil {
		panic(ErrUnexpectedMessageFormat)
	}

	return matched
}
