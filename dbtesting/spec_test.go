package dbtesting

import (
	"github.com/ghthor/gospec"
	"testing"
)

func TestUnitSpecs(t *testing.T) {
	r := gospec.NewRunner()

	r.AddSpec(DescribeSpecifyExecutor)

	gospec.MainGoTest(r, t)
}

func TestIntegrationSpecs(t *testing.T) {
	r := gospec.NewRunner()

	gospec.MainGoTest(r, t)
}
