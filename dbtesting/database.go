package dbtesting

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
)

type TestDatabase struct {
	mysql.Conn
	name   string
	schema string
}

func (t *TestDatabase) Create() error {
	_, _, err := t.Conn.Query("CREATE DATABASE `%s` DEFAULT COLLATE = 'utf8_general_ci'", t.name)
	if err != nil {
		return err
	}
	return t.Use(t.name)
}

func (t *TestDatabase) Drop() error {
	_, _, err := t.Conn.Query("drop database `%s`", t.name)
	return err
}

func (t *TestDatabase) SetSchema(schema string) error {
	if t.schema != "" {
		return errors.New("schema already set")
	}

	res, err := t.Start(string(schema))
	if err != nil {
		return err
	}

	// Must read all results to unlock mysql.Conn for next queries
	for res.MoreResults() {
		if res, err = res.NextResult(); err != nil {
			return err
		}
	}

	t.schema = schema
	return nil
}

func genSuffix() (string, error) {
	suffix := make([]byte, 16)
	n, err := rand.Read(suffix)
	if n != len(suffix) || err != nil {
		return "", err
	}

	return hex.EncodeToString(suffix), nil
}

func newTestDatabase(basename string, c mysql.Conn, genSuffix func() (string, error)) (*TestDatabase, error) {
	suffix, err := genSuffix()
	if err != nil {
		return nil, err
	}

	return &TestDatabase{c, basename + "_" + suffix, ""}, nil
}

func NewTestDatabase(basename string, c mysql.Conn) (*TestDatabase, error) {
	return newTestDatabase(basename, c, genSuffix)
}
