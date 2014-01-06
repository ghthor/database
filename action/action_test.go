package action

import (
	"github.com/ghthor/gospec"
)

type ActionSpec interface {
	Describe(c gospec.Context)
}

func DescribeActions(c gospec.Context) {
	SpecifyActions := func(actions []ActionSpec) {
		for _, a := range actions {
			a.Describe(c)
		}
	}

	SpecifyActions([]ActionSpec{})
}
