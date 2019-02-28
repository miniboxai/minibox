package storage

import "errors"

var ErrMustPrimarykeyValue = errors.New("invalid record, primary value is zero")
var ErrMustStructValue = errors.New("invalid value, must struct value")
