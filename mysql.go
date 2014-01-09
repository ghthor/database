package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
)

type MysqlDatabase struct {
	mysql.Conn
	name   string
	schema string
}

func (t *MysqlDatabase) Create() error {
	_, _, err := t.Conn.Query("CREATE DATABASE `%s` DEFAULT COLLATE = 'utf8_general_ci'", t.name)
	if err != nil {
		return err
	}
	return t.Use(t.name)
}

func (t *MysqlDatabase) Drop() error {
	_, _, err := t.Conn.Query("drop database `%s`", t.name)
	return err
}

func checkIfDatabaseExists(c mysql.Conn, db string) (bool, error) {
	row, _, err := c.QueryFirst("select schema_name from information_schema.schemata where schema_name = '%s'", db)
	if err != nil {
		return false, err
	}

	return len(row) != 0, nil
}

func (t *MysqlDatabase) Exists() (bool, error) {
	return checkIfDatabaseExists(t, t.name)
}

func (t *MysqlDatabase) SetSchema(schema string) error {
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

func newMysqlDatabase(basename string, c mysql.Conn, genSuffix func() (string, error)) (*MysqlDatabase, error) {
	suffix, err := genSuffix()
	if err != nil {
		return nil, err
	}

	return &MysqlDatabase{c, basename + "_" + suffix, ""}, nil
}

func NewMysqlDatabase(basename string, c mysql.Conn) (*MysqlDatabase, error) {
	return newMysqlDatabase(basename, c, genSuffix)
}
