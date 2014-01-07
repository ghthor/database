package config

import (
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	"testing"
)

func TestUnitSpecs(t *testing.T) {
	r := gospec.NewRunner()

	r.AddSpec(DescribeConfigLoading)

	gospec.MainGoTest(r, t)
}

func DescribeConfigLoading(c gospec.Context) {
	c.Specify("Config can be parsed from json file", func() {
		expectedConfig := Config{
			"dbuser",
			"dbpassword",
			"dbname",
			"filepath/to/filedb",
		}

		config, err := ReadFromFile("config.example.json")
		c.Assume(err, IsNil)
		c.Expect(config, Equals, expectedConfig)
	})
}
