package database

import (
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
)

func DescribeUpdateStmtResult(c gospec.Context) {
	c.Specify("an update statement's result", func() {
		updateResult := &UpdateResult{&MockResult{"(Rows matched: 1  Changed: 0  Warnings: 0"}}

		c.Specify("can identify the number of rows that were matched", func() {
			c.Expect(updateResult.MatchedRows(), Equals, uint64(1))

			c.Specify("and panics if the message is in an unexpected format", func() {
				updateResult.Result = &MockResult{"unexpected format: 0 panic: 0 mode: 0"}

				defer func() {
					e := recover()
					c.Expect(e, Not(IsNil))
					c.Expect(e, Equals, ErrUnexpectedMessageFormat)
				}()
				updateResult.MatchedRows()
			})
		})
	})
}
