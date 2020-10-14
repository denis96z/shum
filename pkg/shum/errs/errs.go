package errs

import (
	"github.com/pkg/errors"
)

var (
	New    = errors.New
	Errorf = errors.Errorf
	Wrap   = errors.Wrap
	Wrapf  = errors.Wrapf
	Cause  = errors.Cause
)
