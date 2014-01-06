package action

import (
	"github.com/ghthor/gospec"
	"testing"
)

func TestUnitSpecs(t *testing.T) {
	r := gospec.NewRunner()

	r.AddSpec(DescribeDatatypeConversions)
	r.AddSpec(DescribeActions)

	gospec.MainGoTest(r, t)
}
