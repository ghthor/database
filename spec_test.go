package database

import (
	"github.com/ghthor/gospec"
	"testing"
)

func TestUnitSpecs(t *testing.T) {
	r := gospec.NewRunner()

	r.AddSpec(DescribeUpdateStmtResult)
	r.AddSpec(DescribeMockStmt)

	r.AddSpec(DescribeExecutorRegistry)

	gospec.MainGoTest(r, t)
}

func TestIntegrationSpecs(t *testing.T) {
	r := gospec.NewRunner()

	r.AddSpec(DescribeDatabaseIntegration)

	gospec.MainGoTest(r, t)
}
