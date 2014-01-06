package action

import (
	"github.com/ghthor/database/datatype"
	"strconv"
	"strings"
)

// Generics
type (
	Float struct {
		Str    string
		Native float64
	}

	TimeSpan struct {
		Str   string
		Hours uint64
		Mins  uint64
	}
)

// Ids
type (
	Id struct {
		Str    string
		Native datatype.Id
	}
)

func (f *Float) Parse() error {
	val, err := strconv.ParseFloat(f.Str, 64)
	if err != nil {
		return ErrInvalidFloat
	}
	f.Native = val
	return nil
}

func (t *TimeSpan) Parse() error {
	parts := strings.Split(t.Str, ":")
	if len(parts) != 2 {
		return ErrInvalidTimeSpan
	}

	hours, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return ErrInvalidTimeSpan
	}

	mins, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return ErrInvalidTimeSpan
	}

	t.Hours = hours
	t.Mins = mins
	return nil
}

func (p *Id) Parse() error {
	id, err := strconv.ParseUint(p.Str, 10, 64)
	if err != nil {
		return ErrInvalidId
	}
	p.Native = datatype.Id(id)
	return nil
}
