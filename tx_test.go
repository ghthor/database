package database

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/ghthor/database/datatype"
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"github.com/ziutek/mymysql/mysql"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type MockMysqlTx struct {
	CommitWasCalled   bool
	RollbackWasCalled bool

	DoWasCalled bool
	DoFunc      func(mysql.Stmt) mysql.Stmt
}

func (t *MockMysqlTx) Commit() error {
	t.CommitWasCalled = true
	return nil
}

func (t *MockMysqlTx) Rollback() error {
	t.RollbackWasCalled = true
	return nil
}

func (t *MockMysqlTx) Do(s mysql.Stmt) mysql.Stmt {
	t.DoWasCalled = true
	if t.DoFunc != nil {
		return t.DoFunc(s)
	}
	return s
}

func DescribeTransaction(c gospec.Context) {
	// TODO: Extract filepath stuff into a type that can me mocked so no filesystem calls are made during this test
	tmp, err := ioutil.TempDir("", "transaction-spec")
	c.Assume(err, IsNil)

	defer func() { c.Assume(os.RemoveAll(tmp), IsNil) }()

	tx := newTransaction(&MockMysqlTx{}, tmp)

	type TestFile struct {
		file     datatype.FormFile
		bytes    []byte
		sha1name string
	}

	filenames := map[string]string{
		"png": "dbtesting/image_test.png",
		"txt": "dbtesting/text_test.txt",
	}
	files := make(map[string]*TestFile, len(filenames))

	for k, v := range filenames {
		file, err := testFile(v)
		c.Assume(err, IsNil)

		fileBytes, err := ioutil.ReadFile(v)
		c.Assume(err, IsNil)

		h := sha1.New()
		_, err = io.Copy(h, bytes.NewReader(fileBytes))
		c.Assume(err, IsNil)

		sha1Name := hex.EncodeToString(h.Sum(nil)) + filepath.Ext(v)

		files[k] = &TestFile{file, fileBytes, sha1Name}
	}

	c.Specify("a transaction", func() {
		filename, err := tx.SaveFile(files["png"].file)
		c.Assume(err, IsNil)
		c.Assume(filename, Equals, files["png"].sha1name)

		// TODO: Specify RollbackError behavior
		c.Specify("will rollback", func() {
			c.Specify("during a failed mysql statment", func() {
				stmt := &MockStmt{
					RunFunc: func(...interface{}) (mysql.Result, error) {
						return nil, errors.New("run failed")
					},
				}

				_, err := tx.Run(stmt)
				c.Assume(err, Not(IsNil))
				c.Assume(err.Error(), Equals, "run failed")

				c.Expect(tx.tx.(*MockMysqlTx).RollbackWasCalled, IsTrue)

				c.Specify("and will remove any successfully saved files", func() {
					_, err := os.Stat(path.Join(tmp, files["png"].sha1name))
					c.Expect(os.IsNotExist(err), IsTrue)
				})
			})

			c.Specify("during a failed save file action", func() {
				_, err := tx.saveFile(files["txt"].file, func(datatype.FormFile, string) (string, error) {
					return "", errors.New("error saving file")
				})
				c.Assume(err, Not(IsNil))
				c.Assume(err.Error(), Equals, "error saving file")

				c.Expect(tx.tx.(*MockMysqlTx).RollbackWasCalled, IsTrue)

				_, err = os.Stat(path.Join(tmp, files["txt"].sha1name))
				c.Expect(os.IsNotExist(err), IsTrue)

				c.Specify("and will remove any successfully saved files", func() {
					_, err := os.Stat(path.Join(tmp, files["png"].sha1name))
					c.Expect(os.IsNotExist(err), IsTrue)
				})
			})
		})
	})
}
