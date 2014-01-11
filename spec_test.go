package database

import (
	"github.com/ghthor/database/config"
	"github.com/ghthor/gospec"
	"log"
	"testing"
)

var cfg config.Config

func init() {
	var err error
	cfg, err = config.ReadFromFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
}

func TestUnitSpecs(t *testing.T) {
	r := gospec.NewRunner()

	r.AddSpec(DescribeUpdateStmtResult)
	r.AddSpec(DescribeMockMysqlConn)
	r.AddSpec(DescribeMockStmt)

	r.AddSpec(DescribeTransaction)

	r.AddSpec(DescribeExecutorRegistry)

	gospec.MainGoTest(r, t)
}

func TestIntegrationSpecs(t *testing.T) {
	r := gospec.NewRunner()

	r.AddSpec(DescribeDatabaseIntegration)
	r.AddSpec(DescribeMysqlDatabaseIntegration)

	gospec.MainGoTest(r, t)
}
