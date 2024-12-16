package util

import (
	"errors"

	"github.com/peterbourgon/ff/v4"
)

var ErrCliRequiredFlags = errors.New(ff.ErrHelp.Error())
