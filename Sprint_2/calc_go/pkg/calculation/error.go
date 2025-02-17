package calculation

import (
	"errors"
)

var (
	EmptyExpressionErr     = errors.New(`empty expression`)
	IncorrectExpressionErr = errors.New(`incorrect expression`)
	NumToPopMErr           = errors.New(`numtopop > length of slice of nums`)
	NumToPopZeroErr        = errors.New(`numtopop <= 0`)
	NthToPopErr            = errors.New(`no operator to pop`)
	DvsByZeroErr           = errors.New(`division by zero`)
)
