package progressbar

import "errors"

var (
	ErrInvalidPbWidth = errors.New("invalid progress bar width")
	ErrInvalidPbOut   = errors.New("invalid destination type")
)
