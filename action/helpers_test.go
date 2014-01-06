package action

import (
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
)

func DescribeDatatypeConversions(c gospec.Context) {
	c.Specify("a time span", func() {
		span := TimeSpan{
			Str: "10:30",
		}

		c.Assume(span.Parse(), IsNil)
		c.Assume(span.Hours, Equals, uint64(10))
		c.Assume(span.Mins, Equals, uint64(30))
	})
}
