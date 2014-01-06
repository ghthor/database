package action

import "errors"

var (
	ErrInvalidId       = errors.New("invalid id")
	ErrInvalidFilename = errors.New("invalid filename")
	ErrInvalidFloat    = errors.New("invalid float")
	ErrInvalidTimeSpan = errors.New("invalid timeSpan")
)
