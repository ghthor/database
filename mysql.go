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
	Name   string
	schema string
}

func (t *MysqlDatabase) Create() error {
	_, _, err := t.Conn.Query("CREATE DATABASE `%s` DEFAULT COLLATE = 'utf8_general_ci'", t.Name)
	if err != nil {
		return err
	}
	return t.Use(t.Name)
}

func (t *MysqlDatabase) Drop() error {
	_, _, err := t.Conn.Query("drop database `%s`", t.Name)
	return err
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
